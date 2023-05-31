package db

import (
	"context"
	pgxzap "github.com/jackc/pgx-zap"
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
	tlogger := pgxzap.NewLogger(logger.DesugarZap)
	tracer := &tracelog.TraceLog{
		Logger:   tlogger,
		LogLevel: tracelog.LogLevelTrace,
	}

	cfg, err := pgxpool.ParseConfig(config.Database.DSN())
	cfg.ConnConfig.Tracer = tracer

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	dbpool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		logger.Zap.Fatalf("Unable to create connection pool: %v", err)
	}

	if err = dbpool.Ping(ctx); err != nil {
		logger.Zap.Fatalf("Unable to pinging connection pool: %v", err)
	}

	return Database{
		Pool: dbpool,
	}
}
