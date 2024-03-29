package db

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	pgxzero "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"

	"github.com/kirychukyurii/wasker/internal/config"
	"github.com/kirychukyurii/wasker/internal/lib"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
	"github.com/kirychukyurii/wasker/internal/pkg/server/interceptor/requestid"
)

type Database struct {
	Pool *pgxpool.Pool
}

func New(config config.Config, logger log.Logger) Database {
	ctx := context.Background()

	contextFunc := func(ctx context.Context, zeroCtx zerolog.Context) zerolog.Context {
		reqId := lib.FromContext(ctx, requestid.XRequestIDCtxKey{})
		if reqId != nil {
			return zeroCtx.Str(requestid.XRequestIDKey, reqId.(string))
		}

		return zeroCtx
	}

	// TODO: log messages with level trace (now it info)
	dblogger := pgxzero.NewLogger(logger.Log, pgxzero.WithoutPGXModule(), pgxzero.WithContextFunc(contextFunc))
	dblevel, err := tracelog.LogLevelFromString(logger.Log.GetLevel().String())
	if err != nil {
		logger.Log.Err(err).Msg("setup a pgx tracing level to default: info")
		dblevel = tracelog.LogLevelInfo
	}

	cfg, err := pgxpool.ParseConfig(config.Database.DSN())
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("parsing database connection string")
	}

	cfg.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   dblogger,
		LogLevel: dblevel,
	}

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	dbpool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("create connection pool")
	}

	conn, err := dbpool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("acquire connection from pool")
	}

	if err = conn.Ping(ctx); err != nil {
		logger.Log.Fatal().Err(err).Msg("pinging connection pool")
	}

	return Database{
		Pool: dbpool,
	}
}

func (a *Database) Dialect() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
