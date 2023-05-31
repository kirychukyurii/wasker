package serve

import (
	"context"
	"github.com/kirychukyurii/wasker/internal/config"
	"github.com/kirychukyurii/wasker/internal/pkg"
	"github.com/kirychukyurii/wasker/internal/pkg/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func init() {
	pf := Command.PersistentFlags()

	pf.StringVarP(&config.Path, "config", "c", "config/config.yml",
		"this parameter is used to start the service application")
}

var Command = &cobra.Command{
	Use:          "serve",
	Short:        "",
	Example:      "",
	SilenceUsage: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.SetConfigFile(config.Path)
		if err := viper.ReadInConfig(); err != nil {
			panic(errors.Wrap(err, "failed to read config"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		app := fx.New(
			Module,
			fx.WithLogger(
				func(logger logger.Logger) fxevent.Logger {
					var l fxevent.ZapLogger

					l.UseLogLevel(zap.DebugLevel)
					l.Logger = logger.DesugarZap

					return &l
				},
			),
		)

		app.Run()
	},
}

var Module = fx.Options(
	config.Module,
	pkg.Module,
	fx.Invoke(runApplication),
)

func runApplication(lifecycle fx.Lifecycle, logger logger.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Zap.Info("Starting application")

			return nil
		},
		OnStop: func(context.Context) error {
			logger.Zap.Info("Stopping application")

			return nil
		},
	})
}
