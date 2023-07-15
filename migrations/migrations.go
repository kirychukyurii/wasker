package migrations

import (
	"context"
	"github.com/kirychukyurii/wasker/internal/app/directory/model"
	"github.com/kirychukyurii/wasker/internal/app/directory/repository"
	model2 "github.com/kirychukyurii/wasker/internal/model"
	"github.com/kirychukyurii/wasker/migrations/schema"

	"github.com/jackc/tern/v2/migrate"
	"go.uber.org/fx"

	"github.com/kirychukyurii/wasker/internal/pkg/db"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
)

// Migrate applies all migrations to the database.
func Migrate(ctx context.Context, shutdowner fx.Shutdowner, logger log.Logger, db db.Database, scopeRepository repository.ScopeRepository) {
	conn, err := db.Pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("acquire connection from pool")
	}

	m, err := migrate.NewMigratorEx(context.Background(), conn.Conn(), "public.schema_version", &migrate.MigratorOptions{})
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("create new Migrator")
	}

	if err := m.LoadMigrations(schema.EmbedMigrations); err != nil {
		logger.Log.Fatal().Err(err).Msg("load migrations")
	}

	if err := m.Migrate(ctx); err != nil {
		logger.Log.Fatal().Err(err).Msg("start migrations")
	}

	for _, scope := range model.DefaultPermission {
		var scopeId int64
		l := logger.Log.With().Str("scope", scope.Scope).Logger()

		scopeResult, err := scopeRepository.Query(ctx, &model.ScopeQueryParam{Query: model2.QueryParam{Name: scope.Scope}})
		if err != nil {
			l.Fatal().Err(err).Msg("query row")
		}

		if v := len(scopeResult.List); v == 0 {
			l.Info().Msg("inserting permission scope")
			var s model.Scope
			s.Name = scope.Scope

			scopeId, err = scopeRepository.Create(ctx, &s)
			if err != nil {
				l.Fatal().Err(err).Msg("query row")
			}
		} else {
			scopeId = scopeResult.List[0].Id
			l.Info().Msg("permission scope already exists")
		}

		for _, endpoint := range scope.Endpoint {
			l := logger.Log.With().Str("scope", scope.Scope).Str("scope-endpoint", endpoint.Name).Logger()

			endpointResult, err := scopeRepository.QueryEndpoint(ctx, &model.ScopeEndpointQueryParam{Query: model2.QueryParam{Name: endpoint.Name}})
			if err != nil {
				l.Fatal().Err(err).Msg("query row")
			}

			if v := len(endpointResult.List); v == 0 {
				l.Info().Msg("inserting permission scope endpoint")

				var e model.ScopeEndpoint
				e.Name = endpoint.Name
				e.Bit = endpoint.Bit
				e.Scope.Id = scopeId

				_, err = scopeRepository.CreateEndpoint(ctx, &e)
				if err != nil {
					l.Fatal().Err(err).Msg("query row")
				}
			} else {
				l.Info().Msg("permission scope endpoint already exists")
			}
		}
	}

	if err := shutdowner.Shutdown(fx.ExitCode(0)); err != nil {
		logger.Log.Error().Err(err).Msg("shutdown application")
	}
}
