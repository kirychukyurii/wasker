package db

import (
	"context"
	pgxzero "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/kirychukyurii/wasker/internal/config"
	"github.com/kirychukyurii/wasker/internal/pkg/logger"
)

type Database struct {
	Pool *pgxpool.Pool
}

func New(config config.Config, logger logger.Logger) Database {
	ctx := context.Background()

	// TODO: log messages with level trace (now it info)
	tlogger := pgxzero.NewLogger(logger.Log)
	tracer := &tracelog.TraceLog{
		Logger:   tlogger,
		LogLevel: tracelog.LogLevelTrace,
	}

	cfg, err := pgxpool.ParseConfig(config.Database.DSN())
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed when parsing database connection string")
	}

	cfg.ConnConfig.Tracer = tracer

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	dbpool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to create connection pool")
	}

	if err = dbpool.Ping(ctx); err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to pinging connection pool")
	}

	return Database{
		Pool: dbpool,
	}
}
