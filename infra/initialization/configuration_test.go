package initialization

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const pathConfig = ".config"

func TestExtractConfiguration(t *testing.T) {
	configRaw := GetConfigurationRaw(pathConfig)
	require.NotNil(t, configRaw)

	configuration, errParse := NewConfiguration(configRaw)
	require.NoError(t, errParse)
	require.NotNil(t, configuration)
	require.NotZero(t, configuration.PortRPC)
	require.NotZero(t, configuration.PortHTTP)
	require.NotZero(t, configuration.NameLogfile)
	require.NotZero(t, configuration.Sites)
	require.NotZero(t, configuration.Sites[0].Name)

	require.NotZero(t, configuration.Sites[0].OffsetUTCHours)
	require.NotZero(t, configuration.Sites[0].TimestampDSTSpring)
	require.NotZero(t, configuration.Sites[0].TimestampDSTWinter)
}
