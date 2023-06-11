package db

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	pgxzero "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"

	"github.com/kirychukyurii/wasker/internal/config"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
)

type Database struct {
	Pool *pgxpool.Pool
}

func New(config config.Config, logger log.Logger) Database {
	ctx := context.Background()

	// TODO: log messages with level trace (now it info)
	tlogger := pgxzero.NewLogger(logger.Log)
	tracer := &tracelog.TraceLog{
		Logger:   tlogger,
		LogLevel: tracelog.LogLevelTrace,
	}

	cfg, err := pgxpool.ParseConfig(config.Database.DSN())
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("parsing database connection string")
	}

	cfg.ConnConfig.Tracer = tracer

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

// ColumnDataPair describes a piece of data that is stored in a database table column
type ColumnDataPair struct {
	Column string
	Data   interface{}
}
