package serve

import (
	"context"
	"github.com/kirychukyurii/wasker/internal/config"
	"github.com/kirychukyurii/wasker/internal/controller"
	"github.com/kirychukyurii/wasker/internal/pkg"
	"github.com/kirychukyurii/wasker/internal/pkg/db"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
	"github.com/kirychukyurii/wasker/internal/repository"
	"github.com/kirychukyurii/wasker/internal/server"
	"github.com/kirychukyurii/wasker/internal/service"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"net"
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
				func(logger log.Logger) fxevent.Logger {
					return &log.FxLogger{
						Logger: &logger.Log,
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
	server.Module,
	repository.Module,
	service.Module,
	controller.Module,
	fx.Invoke(runApplication),
)

func runApplication(lifecycle fx.Lifecycle, cfg config.Config, logger log.Logger, db db.Database,
	httpServer server.HttpServer, grpcServer server.GrpcServer) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Log.Info().Str("http-listen", cfg.Http.ListenAddr()).Str("grpc-listen", cfg.Grpc.ListenAddr()).Msg("starting application")

			go func() {
				l, err := net.Listen("tcp", cfg.Grpc.ListenAddr())
				if err != nil {
					logger.Log.Fatal().Err(err).Msgf("error in listening on port :", cfg.Grpc.Port)
				}

				// the gRPC server
				if err := grpcServer.Server.Serve(l); err != nil {
					logger.Log.Fatal().Err(err).Msg("unable to start server")
				}
			}()

			go func() {
				l, err := net.Listen("tcp", cfg.Http.ListenAddr())
				if err != nil {
					logger.Log.Fatal().Err(err).Msgf("failed listen :%d", cfg.Http.Port)
				}

				// the HTTP server
				if err = httpServer.Server.Serve(l); err != nil {
					if !errors.Is(err, http.ErrServerClosed) {
						logger.Log.Fatal().Err(err).Msg("Error to Start Application")
					}
				}
			}()

			return nil
		},
		OnStop: func(context.Context) error {
			logger.Log.Info().Msg("stopping application")
			db.Pool.Close()
			if err := httpServer.Server.Close(); err != nil {
				if !errors.Is(err, http.ErrServerClosed) {
					logger.Log.Debug().Err(err).Msg("failed to stop http server")
				}
			}

			grpcServer.Server.GracefulStop()

			return nil
		},
	})
}
