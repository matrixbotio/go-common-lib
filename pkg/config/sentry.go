package config

type SentryConfig struct {
	DSN           string  `envconfig:"SENTRY_DSN"`
	EnableTracing bool    `envconfig:"SENTRY_ENABLE_TRACING" default:"true"`
	SampleRate    float64 `envconfig:"SENTRY_SAMPLE_RATE" default:"0.2"`
}
