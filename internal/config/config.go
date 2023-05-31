package config

import (
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
