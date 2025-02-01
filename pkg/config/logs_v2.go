package config

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

type LogLevel string

const (
	LogLevelInfo    LogLevel = "log"
	LogLevelVerbose LogLevel = "verbose"
	LogLevelWarn    LogLevel = "warn"
	LogLevelError   LogLevel = "error"
	LogLevelFatal   LogLevel = "critical"
)

type LogsConfig struct {
	LogToES              bool     `envconfig:"LOG_TO_ES" default:"true"`
	SamplingInitial      int      `envconfig:"LOG_SAMPLING_INITIAL" default:"300"`
	SamplingThereafter   int      `envconfig:"LOG_SAMPLING_THEREAFTER" default:"100"`
	AppLevelRaw          LogLevel `envconfig:"LOG_LEVEL" default:"log"`
	RMQWorkerLibLevelRaw LogLevel `envconfig:"LOG_LEVEL_RMQWORKERLIB" default:"log"`
}

func getZapLogLevel(val LogLevel) zapcore.Level {
	if val == "" {
		return zapcore.DebugLevel
	}

	switch val {
	case LogLevelInfo:
		return zapcore.InfoLevel
	case LogLevelVerbose:
		return zapcore.DebugLevel
	case LogLevelWarn:
		return zapcore.WarnLevel
	case LogLevelError:
		return zapcore.ErrorLevel
	case LogLevelFatal:
		return zapcore.FatalLevel
	default:
		fmt.Printf("invalid log level: %q\n", val)
		return zapcore.DebugLevel
	}
}

func (l LogsConfig) GetAppLevel() zapcore.Level {
	return getZapLogLevel(l.AppLevelRaw)
}

func (l LogsConfig) GetRMQLogLevel() zapcore.Level {
	return getZapLogLevel(l.RMQWorkerLibLevelRaw)
}
