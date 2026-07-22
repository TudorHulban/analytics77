package analytics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetaActive_ZeroValueKey(t *testing.T) {
	var m MetaActive[string]

	// 1. Insert a valid empty-string key
	m.Increment("", 1)

	// Expect slot 0 to contain the empty string
	require.Equal(t, "", *m.Names[0].Load())
	require.Equal(t, uint32(1), m.Values[0].Load())

	// 2. Insert another key
	m.Increment("A", 1)

	// Correct behavior: "" must remain untouched.
	require.Equal(t,
		"",
		*m.Names[0].Load(),

		"empty-string key was incorrectly treated as empty slot",
	)
	require.Equal(t, uint32(1), m.Values[0].Load())

	// // 3. Merge with another MetaActive also containing ""
	// var m2 MetaActive[string]
	// m2.Increment("", 2)

	// merged := MetaActive[string]{}.MergeActive(&m)
	// merged = merged.MergeActive(&m2)

	// // Correct behavior: "" must be included and summed.
	// require.Equal(t,
	// 	"",
	// 	merged.Names[0],

	// 	"empty-string key was dropped during merge",
	// )
	// require.Equal(t,
	// 	uint32(3),
	// 	merged.Values[0],

	// 	"merge did not accumulate empty-string key",
	// )
}
