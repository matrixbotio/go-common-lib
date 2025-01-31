package sentryx

import (
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/matrixbotio/go-common-lib/pkg/config"
	"go.uber.org/zap"
)

func InitSentry(cfg config.SentryConfig) error {
	if cfg.DSN == "" {
		zap.L().Debug("skip sentry init: dsn is not set")
		return nil
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.DSN,
		TracesSampleRate: cfg.SampleRate,
		EnableTracing:    cfg.EnableTracing,
	})
	if err != nil {
		return fmt.Errorf("init: %w", err)
	}
	return nil
}
