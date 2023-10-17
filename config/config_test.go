package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShouldBeLoadConfig(t *testing.T) {
	config, err := LoadConfig("../")

	require.NoError(t, err)
	require.NotEmpty(t, config.Driver)
	require.NotEmpty(t, config.Source)
	require.NotEmpty(t, config.ServerAddress)
}
