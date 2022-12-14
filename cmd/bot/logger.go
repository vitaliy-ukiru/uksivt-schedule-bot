package main

import (
	"time"

	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newLoggerConfig(cfg *config.Config) zap.Config {
	logger := cfg.Logger
	return zap.Config{
		Level:       logger.Level,
		Development: false,
		Encoding:    logger.Output.Format,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "msg",
			LevelKey:      "level",
			TimeKey:       "ts",
			NameKey:       "logger",
			CallerKey:     "caller",
			StacktraceKey: "stacktrace",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.LowercaseLevelEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format("2006-01-02 15:04:05.000 -07:00"))
			},
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},

		OutputPaths:      logger.Output.Paths,
		ErrorOutputPaths: logger.Output.ErrorPaths,
		InitialFields:    nil,
	}

}
