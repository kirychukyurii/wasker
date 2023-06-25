package migrate

import (
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kirychukyurii/wasker/internal/config"
	repository2 "github.com/kirychukyurii/wasker/internal/directory/repository"
	"github.com/kirychukyurii/wasker/internal/pkg/db"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/kirychukyurii/wasker/internal/pkg"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
	"github.com/kirychukyurii/wasker/migrations"
)

func init() {
	pf := Command.PersistentFlags()

	pf.StringVarP(&config.Path, "config", "c", "config/config.yml",
		"this parameter is used to start the service application")
}

var Command = &cobra.Command{
	Use:          "migrate",
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
	repository2.Module,
	fx.Invoke(runApplication),
)

func runApplication(lifecycle fx.Lifecycle, shutdowner fx.Shutdowner, logger log.Logger, db db.Database, scopeRepository repository2.ScopeRepository) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Log.Info().Msg("starting migrations")

			go migrations.Migrate(ctx, shutdowner, logger, db, scopeRepository)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Log.Info().Msg("stopping application")
			db.Pool.Close()

			return nil
		},
	})
}
