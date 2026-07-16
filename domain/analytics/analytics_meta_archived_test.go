package analytics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetaArchivedMergeActive(t *testing.T) {
	nameRO := "RO"
	nameDE := "DE"
	nameUS := "US"

	dst := MetaArchived[string]{
		Names:  [7]string{nameRO, nameUS},
		Values: [7]uint32{2, 5},
	}

	var active MetaActive[string]

	active.Names[0].Store(&nameRO)
	active.Values[0].Store(3)
	active.setOccupied(0)

	active.Names[1].Store(&nameDE)
	active.Values[1].Store(7)
	active.setOccupied(1)

	archived := dst.MergeActive(&active)

	// RO should be 2 + 3 = 5
	require.Equal(t, uint32(5), archived.Count(nameRO))

	// US should remain 5
	require.Equal(t, uint32(5), archived.Count(nameUS))

	// DE should be 7
	require.Equal(t, uint32(7), archived.Count(nameDE))
}

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

	// 3. Merge with another MetaActive also containing ""
	var m2 MetaActive[string]
	m2.Increment("", 2)

	merged := MetaArchived[string]{}.MergeActive(&m)
	merged = merged.MergeActive(&m2)

	// Correct behavior: "" must be included and summed.
	require.Equal(t,
		"",
		merged.Names[0],

		"empty-string key was dropped during merge",
	)
	require.Equal(t,
		uint32(3),
		merged.Values[0],

		"merge did not accumulate empty-string key",
	)
}

func TestMetaArchived_ZeroValueKey(t *testing.T) {
	// dst contains "" with value 1
	dst := MetaArchived[string]{
		Names:  [7]string{""},
		Values: [7]uint32{1},
	}

	// src contains "" with value 2
	src := MetaArchived[string]{
		Names:  [7]string{""},
		Values: [7]uint32{2},
	}

	merged := dst.MergeArchived(src)

	// We do NOT assume ordering. We search for "" with value 3.
	found := false

	for ix := range 7 {
		if merged.Names[ix] == "" && merged.Values[ix] == 3 {
			found = true

			break
		}
	}

	require.True(t,
		found,
		"merged archived did not accumulate zero-value key correctly",
	)
}

func TestMetaArchived_FullLoadStability(t *testing.T) {
	// dst fully populated with 7 entries
	dst := MetaArchived[string]{
		Names:  [7]string{"A", "B", "C", "D", "E", "F", "G"},
		Values: [7]uint32{10, 20, 30, 40, 50, 60, 70},
	}

	// src fully populated with 7 entries
	src := MetaArchived[string]{
		Names:  [7]string{"H", "I", "J", "K", "L", "M", "N"},
		Values: [7]uint32{15, 25, 35, 45, 55, 65, 75},
	}

	merged := dst.MergeArchived(src)

	// Expected top 7 values across both sets:
	// 75 (N), 70 (G), 65 (M), 60 (F), 55 (L), 50 (E), 45 (K)
	expected := map[string]uint32{
		"N": 75,
		"G": 70,
		"M": 65,
		"F": 60,
		"L": 55,
		"E": 50,
		"K": 45,
	}

	// Validate that all expected keys exist with correct values
	for key, val := range expected {
		require.Equalf(
			t,
			val,
			merged.Count(key),

			"missing or incorrect value for key %q",
			key,
		)
	}

	// Ensure no unexpected keys appear in the top 7
	for ix := range 7 {
		name := merged.Names[ix]
		if name == "" {
			continue
		}

		_, exists := expected[name]

		require.Truef(
			t,
			exists,

			"unexpected key %q found in merged top-7",
			name,
		)
	}
}
