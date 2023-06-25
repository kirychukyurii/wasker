package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var (
	// Path used for setting config path
	Path string

	// Module exports dependency
	Module = fx.Options(
		fx.Provide(New),
	)
)

var DefaultConfig = Config{
	Http: &HttpConfig{
		Host: "0.0.0.0",
		Port: 8080,
	},
	Log: &LogConfig{
		Level:     "debug",
		Directory: "/tmp/tasker-api",
	},
	Auth: &Auth{},
	Database: &DatabaseConfig{
		MaxLifetime:  7200,
		MaxOpenConns: 150,
		MaxIdleConns: 50,
	},
}

type Config struct {
	Http     *HttpConfig     `mapstructure:"http"`
	Grpc     *GrpcConfig     `mapstructure:"grpc"`
	Consul   *ConsulConfig   `mapstructure:"consul"`
	Log      *LogConfig      `mapstructure:"log"`
	Auth     *Auth           `mapstructure:"auth"`
	Database *DatabaseConfig `mapstructure:"database"`
}

type HttpConfig struct {
	Host string `mapstructure:"host" validate:"ipv4"`
	Port int    `mapstructure:"port" validate:"gte=1,lte=65535"`
}

type GrpcConfig struct {
	Timeout int    `mapstructure:"timeout" validate:"gte=1,lte=65535"`
	Host    string `mapstructure:"host" validate:"ipv4"`
	Port    int    `mapstructure:"port" validate:"gte=1,lte=65535"`
}

type ConsulConfig struct {
	Scheme     string `mapstructure:"scheme"`
	Datacenter string `mapstructure:"datacenter"`
	Host       string `mapstructure:"host" validate:"ipv4"`
	Port       int    `mapstructure:"port" validate:"gte=1,lte=65535"`
}

type LogConfig struct {
	Level     string `mapstructure:"level"`     // debug, info, warn, error, dpanic, panic, fatal
	Format    string `mapstructure:"format"`    // json, console
	Directory string `mapstructure:"directory"` // log storage path
}

type Auth struct {
	TokenExpired int    `mapstructure:"token_ttl"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
}

type DatabaseConfig struct {
	Name     string `mapstructure:"name"`
	Host     string `mapstructure:"host" validate:"ipv4"`
	Port     int    `mapstructure:"port" validate:"gte=1,lte=65535"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`

	MaxLifetime  int `mapstructure:"max_connection_lifetime"`
	MaxOpenConns int `mapstructure:"max_opened_connections"`
	MaxIdleConns int `mapstructure:"max_idle_connections"`
}

func New() Config {
	config := DefaultConfig

	if err := viper.Unmarshal(&config); err != nil {
		panic(errors.Wrap(err, "failed to marshal config"))
	}

	return config
}

func (a *DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?pool_max_conns=10",
		a.Username, a.Password, a.Host, a.Port, a.Name,
	)
}

func (a *HttpConfig) ListenAddr() string {
	if err := validator.New().Struct(a); err != nil {
		return fmt.Sprintf("%s:%d", DefaultConfig.Http.Host, DefaultConfig.Http.Port)
	}

	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

func (a *GrpcConfig) ListenAddr() string {
	if err := validator.New().Struct(a); err != nil {
		return fmt.Sprintf("%s:%d", DefaultConfig.Grpc.Host, DefaultConfig.Grpc.Port)
	}

	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

func (a *ConsulConfig) Addr() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}
