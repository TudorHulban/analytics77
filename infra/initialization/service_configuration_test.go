package initialization

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfiguration(t *testing.T) {
	path := "../../cmd/.config"

	config := Configuration(path)
	require.NotNil(t, config)

	fmt.Println(config)
}
