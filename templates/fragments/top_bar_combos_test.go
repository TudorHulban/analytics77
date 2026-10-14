package fragments

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestComboSelectMonth(t *testing.T) {
	combo := NewComboSelectMonth(
		time.Now(),
		3,
		3,
	)
	require.NotZero(t, combo)

	node := combo.Build()
	require.NotNil(t, node)
	require.NotZero(t, node.Canonical())
}
