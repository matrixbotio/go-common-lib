package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matrixbotio/go-common-lib/pkg/config"
)

func TestGetAppConfig_Set(t *testing.T) {
	os.Clearenv()

	err := os.Setenv("ENV", "wow")
	require.NoError(t, err)

	c, err := config.GetAppConfig()

	assert.NoError(t, err)
	assert.Equal(t, "wow", c.Env)
}

func TestGetAppConfig_NotSet(t *testing.T) {
	os.Clearenv()

	c, err := config.GetAppConfig()

	assert.NoError(t, err)
	assert.Equal(t, "none", c.Env)
}

func TestGetAppConfig_SetToEmpty(t *testing.T) {
	os.Clearenv()

	err := os.Setenv("ENV", "")
	require.NoError(t, err)

	_, err = config.GetAppConfig()

	assert.ErrorIs(t, err, config.EmptyFieldError{"env"})
}
