package analytics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetaActive_ZeroValueKey(t *testing.T) {
	var m1 MetaActive[string]

	// 1. Insert a valid empty-string key
	m1.Increment("", 1)

	// Expect slot 0 to contain the empty string
	require.Equal(t,
		"",
		*m1.Names[0].Load(),
	)
	require.Equal(t,
		uint32(1),
		m1.Values[0].Load(),
	)

	// 2. Insert another key
	m1.Increment("A", 1)

	// Correct behavior: "" must remain untouched.
	require.Equal(t,
		"",
		*m1.Names[0].Load(),

		"empty-string key was incorrectly treated as empty slot",
	)
	require.Equal(t,
		uint32(1),
		m1.Values[0].Load(),
	)

	// 3. Merge with another MetaActive also containing ""
	var m2 MetaActive[string]

	m2.Increment("", 2)

	m1.MergeFrom(&m2)

	// Correct behavior: "" must be included and summed.
	require.Equal(t,
		"",
		*m1.Names[0].Load(),

		"empty-string key was dropped during merge",
	)
	require.Equal(t,
		uint32(3),
		m1.Values[0].Load(),

		"merge did not accumulate empty-string key",
	)

	var m3 MetaActive[string]

	m3.MergeFrom(&m2)

	require.Equal(t,
		"",
		*m3.Names[0].Load(),

		"empty-string key was dropped during merge",
	)
}
