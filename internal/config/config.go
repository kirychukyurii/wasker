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
