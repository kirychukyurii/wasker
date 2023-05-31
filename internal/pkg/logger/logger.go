package logger

import (
	"github.com/kirychukyurii/wasker/internal/config"
	"github.com/kirychukyurii/wasker/internal/constants"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

type Logger struct {
	Zap        *zap.SugaredLogger
	DesugarZap *zap.Logger
}

func New(config config.Config) Logger {
	var options []zap.Option

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(constants.TimeFormat)

	level := zap.NewAtomicLevelAt(toLevel(config.Log.Level))

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		toWriter(),
		level,
	)

	stackLevel := zap.NewAtomicLevel()
	stackLevel.SetLevel(zap.WarnLevel)
	options = append(options,
		zap.AddCaller(),
		zap.AddStacktrace(stackLevel),
	)

	logger := zap.New(core, options...)

	return Logger{
		Zap:        logger.Sugar(),
		DesugarZap: logger,
	}
}

func toLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "dpanic":
		return zap.DPanicLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func toWriter() zapcore.WriteSyncer {
	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
	)
}
