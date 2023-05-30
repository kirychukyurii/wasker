package logger

import (
	"github.com/kirychukyurii/wasker/internal/constants"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Logger struct {
	Zap        *zap.SugaredLogger
	DesugarZap *zap.Logger
}

func New() Logger {
	var options []zap.Option

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(constants.TimeFormat)

	level := zap.NewAtomicLevelAt(zap.DebugLevel)

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

func toWriter() zapcore.WriteSyncer {
	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
	)
}
