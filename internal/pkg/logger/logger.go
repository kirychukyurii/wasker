package logger

import (
	"github.com/kirychukyurii/wasker/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"io"
	"os"
	"strings"
	"time"
)

type Logger struct {
	Log zerolog.Logger
}

func New(config config.Config) Logger {
	var log zerolog.Logger
	var output io.Writer = zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339Nano

	log = zerolog.New(output).
		Level(toLevel(config.Log.Level)).
		With().
		Timestamp().
		Caller().
		Logger()

	return Logger{
		Log: log,
	}
}

func toLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "panic":
		return zerolog.PanicLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}
