package sstorage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceStorage(t *testing.T) {
	serviceStorage, errCrServiceStorage := NewServiceStorage(t.TempDir())
	require.NoError(t, errCrServiceStorage)
	require.NotNil(t, serviceStorage)

	t.Cleanup(
		func() {
			require.NoError(t, serviceStorage.Close())
		},
	)
}
