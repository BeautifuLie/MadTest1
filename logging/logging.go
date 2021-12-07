package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitZapLog() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := config.Build()
	logger.Info("Starting programm")
	return logger
}
