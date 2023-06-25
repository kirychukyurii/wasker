package directory

import (
	"context"
	"fmt"
	"net"

	"github.com/google/uuid"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/kirychukyurii/wasker/internal/config"
	"github.com/kirychukyurii/wasker/internal/constants"
	"github.com/kirychukyurii/wasker/internal/directory/controller"
	"github.com/kirychukyurii/wasker/internal/directory/repository"
	"github.com/kirychukyurii/wasker/internal/directory/server/register"
	"github.com/kirychukyurii/wasker/internal/directory/service"
	"github.com/kirychukyurii/wasker/internal/errors"
	"github.com/kirychukyurii/wasker/internal/pkg"
	"github.com/kirychukyurii/wasker/internal/pkg/consul"
	"github.com/kirychukyurii/wasker/internal/pkg/db"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
	"github.com/kirychukyurii/wasker/internal/pkg/server"
)

func init() {
	pf := Command.PersistentFlags()

	pf.StringVarP(&config.Path, "config", "c", "config/config.yml",
		"this parameter is used to start the service application")
}

var Command = &cobra.Command{
	Use:          "directory",
	Short:        "",
	Example:      "",
	SilenceUsage: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.SetConfigFile(config.Path)
		if err := viper.ReadInConfig(); err != nil {
			panic(errors.Wrap(err, "read config"))
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
	repository.Module,
	service.Module,
	controller.Module,
	fx.Invoke(runApplication),
)

func runApplication(lifecycle fx.Lifecycle, cfg config.Config, logger log.Logger, db db.Database,
	grpcServer server.GrpcServer, discovery consul.ServiceDiscovery, controller controller.Controllers) {
	serviceId := fmt.Sprintf("%s-%s", constants.DirectoryServiceName, uuid.New().String())
	subLogger := logger.Log.With().Str("service-id", serviceId).Logger()

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			subLogger.Info().Str("grpc-listen", cfg.Grpc.ListenAddr()).Msg("starting application")

			registration := &consulapi.AgentServiceRegistration{
				ID:      serviceId,
				Name:    constants.DirectoryServiceName,
				Port:    cfg.Grpc.Port,
				Address: cfg.Grpc.Host,
			}

			err := discovery.Client.Agent().ServiceRegister(registration)
			if err != nil {
				subLogger.Fatal().Err(err).Msg("register service")
			}

			go func() {
				l, err := net.Listen("tcp", cfg.Grpc.ListenAddr())
				if err != nil {
					subLogger.Fatal().Err(err).Msgf("listening on port :%d", cfg.Grpc.Port)
				}

				// register controllers on the gRPC server
				register.GrpcDirectoryServiceServers(grpcServer.Server, controller)

				// the gRPC server
				if err := grpcServer.Server.Serve(l); err != nil {
					subLogger.Fatal().Err(err).Msg("start server")
				}
			}()

			return nil
		},
		OnStop: func(context.Context) error {
			subLogger.Info().Msg("stopping application")

			db.Pool.Close()
			grpcServer.Server.GracefulStop()
			err := discovery.Client.Agent().ServiceDeregister(serviceId)
			if err != nil {
				subLogger.Error().Err(err).Msg("deregister service")
			}

			return nil
		},
	})
}
