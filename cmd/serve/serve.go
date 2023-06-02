package serve

import (
	"context"
	"github.com/kirychukyurii/wasker/internal/api/middleware"
	"github.com/kirychukyurii/wasker/internal/config"
	"github.com/kirychukyurii/wasker/internal/pkg"
	"github.com/kirychukyurii/wasker/internal/pkg/db"
	"github.com/kirychukyurii/wasker/internal/pkg/handler"
	"github.com/kirychukyurii/wasker/internal/pkg/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"net/http"
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
				func(log logger.Logger) fxevent.Logger {
					return &logger.FxLogger{
						Logger: &log.Log,
					}
				},
			),
		)

		app.Run()
	},
}

var Module = fx.Options(
	config.Module,
	pkg.Module,
	middleware.Module,
	fx.Invoke(runApplication),
)

func runApplication(lifecycle fx.Lifecycle, cfg config.Config, logger logger.Logger, db db.Database,
	handler handler.HttpHandler, middleware middleware.Middlewares) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Log.Info().Str("http-listen", cfg.Http.ListenAddr()).Msg("Starting application")

			go func() {
				middleware.Setup()

				if err := handler.Engine.Start(cfg.Http.ListenAddr()); err != nil {
					if errors.Is(err, http.ErrServerClosed) {
						logger.Log.Debug().Err(err).Msg("Shutting down the Application")
					} else {
						logger.Log.Fatal().Err(err).Msg("Error to Start Application")
					}
				}
			}()

			return nil
		},
		OnStop: func(context.Context) error {
			logger.Log.Info().Msg("Stopping application")
			db.Pool.Close()

			return nil
		},
	})
}
