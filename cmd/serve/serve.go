package serve

import (
	"context"
	"github.com/kirychukyurii/wasker/internal/pkg"
	"github.com/kirychukyurii/wasker/internal/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func init() {
	_ = Command.PersistentFlags()
}

var Command = &cobra.Command{
	Use:          "serve",
	Short:        "",
	Example:      "",
	SilenceUsage: true,
	PreRun:       func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		f := fx.New(
			Module,
			fx.WithLogger(
				func(log logger.Logger) fxevent.Logger {
					var l fxevent.ZapLogger

					l.UseLogLevel(zap.DebugLevel)
					l.Logger = log.DesugarZap

					return &l
				},
			),
		)

		f.Run()
	},
}

var Module = fx.Options(
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
