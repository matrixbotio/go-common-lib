package sentryx

import (
	"testing"

	"github.com/matrixbotio/go-common-lib/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestInitSentryNoDSN(t *testing.T) {
	// given
	cfg := config.SentryConfig{}

	// when
	err := InitSentry(cfg)

	// then
	require.NoError(t, err)
}

func TestInitSentryEmptyUserNameError(t *testing.T) {
	// given
	cfg := config.SentryConfig{
		DSN: "https://sentry.io",
	}

	// when
	err := InitSentry(cfg)

	// then
	require.ErrorContains(t, err, "empty username")
}

func TestInitSentrySuccess(t *testing.T) {
	// given
	cfg := config.SentryConfig{
		DSN: "https://username@code.ingest.de.sentry.io/code2",
	}

	// when
	err := InitSentry(cfg)

	// then
	require.NoError(t, err)
}
