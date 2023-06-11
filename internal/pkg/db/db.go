package db

import (
	"context"
	"database/sql/driver"
	"reflect"

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
		logger.Log.Fatal().Err(err).Msg("Failed when parsing database connection string")
	}

	cfg.ConnConfig.Tracer = tracer

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	dbpool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to create connection pool")
	}

	conn, err := dbpool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to acquire connection from pool")
	}

	if err = conn.Ping(ctx); err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to pinging connection pool")
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

// GetFields returns you an array of ColumnDataPairs which describe
// a database row.
// It uses the db struct tag to get the table column names
func GetFields(s interface{}) ([]ColumnDataPair, error) {
	var row []ColumnDataPair
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		col := field.Tag.Get("db")
		if col == "" {
			col = field.Name
		}

		val, err := driver.DefaultParameterConverter.ConvertValue(v.Field(i).Interface())
		if err != nil {
			return nil, err
		}
		row = append(row, ColumnDataPair{col, val})
	}
	return row, nil
}
