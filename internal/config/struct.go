package config

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
	Log      *LogConfig      `mapstructure:"log"`
	Auth     *Auth           `mapstructure:"auth"`
	Database *DatabaseConfig `mapstructure:"database"`
}

type HttpConfig struct {
	Host string `mapstructure:"host" validate:"ipv4"`
	Port int    `mapstructure:"port" validate:"gte=1,lte=65535"`
}

type GrpcConfig struct {
	Host string `mapstructure:"host" validate:"ipv4"`
	Port int    `mapstructure:"port" validate:"gte=1,lte=65535"`
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
