package zes

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	timeKey       = "timestamp"
	messageKey    = "message"
	stacktraceKey = "stack"
	levelTextKey  = "level_text"
)

func newConsoleConfig() zapcore.EncoderConfig {
	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.TimeKey = timeKey
	cfg.MessageKey = messageKey
	cfg.StacktraceKey = stacktraceKey
	cfg.LevelKey = levelTextKey

	return cfg
}

func newJsonConfig() zapcore.EncoderConfig {
	cfg := zap.NewProductionEncoderConfig()
	cfg.TimeKey = timeKey
	cfg.MessageKey = messageKey
	cfg.StacktraceKey = stacktraceKey
	cfg.LevelKey = levelTextKey
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder

	return cfg
}
