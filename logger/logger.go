package logger

import (
	"log"
	"os"

	"cloud.google.com/go/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func GlobalBuild(isProduction bool) {
	var err error
	if isProduction {
		logger, err = NewStackdriverLogging()
	} else {
		logger, err = NewDevelopmentLogging()
	}

	if err != nil {
		log.Fatal("can't create logging")
	}

	defer logger.Sync()
}

func NewStackdriverLogging() (*zap.Logger, error) {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:      "severity",
			NameKey:       "logger",
			CallerKey:     "caller",
			StacktraceKey: "stacktrace",
			TimeKey:       "time",
			MessageKey:    "message",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeTime:    zapcore.RFC3339NanoTimeEncoder,
			EncodeLevel:   levelEncode,
			EncodeCaller:  zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build(zap.AddCallerSkip(1), zap.AddStacktrace(zap.DPanicLevel))
}

func NewDevelopmentLogging() (*zap.Logger, error) {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development: true,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:       "severity",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			TimeKey:        "timestamp",
			MessageKey:     "message",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build(zap.AddCallerSkip(1), zap.AddStacktrace(zap.DPanicLevel))
}

func levelEncode(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch l {
	case zapcore.DebugLevel:
		enc.AppendString(logging.Debug.String())
	case zapcore.InfoLevel:
		enc.AppendString(logging.Info.String())
	case zapcore.WarnLevel:
		enc.AppendString(logging.Warning.String())
	case zapcore.ErrorLevel:
		enc.AppendString(logging.Error.String())
	case zapcore.DPanicLevel:
		enc.AppendString(logging.Critical.String())
	case zapcore.PanicLevel:
		enc.AppendString(logging.Alert.String())
	case zapcore.FatalLevel:
		enc.AppendString(logging.Emergency.String())
	}
}

func Info(msg string, fields ...zapcore.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zapcore.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zapcore.Field) {
	logger.Error(msg, fields...)
}

func Critical(msg string, fields ...zapcore.Field) {
	logger.DPanic(msg, fields...)
	os.Exit(1)
}
