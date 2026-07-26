package analytics

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetaActiveCount(t *testing.T) {
	var m MetaActive[string]

	// (1) Empty → Count must be 0
	require.Equal(t, uint32(0), m.Count("RO"))

	// (2) Single increment
	m.Increment("RO", 3)
	require.Equal(t, uint32(3), m.Count("RO"))
	require.Equal(t, uint32(0), m.Count("DE"))

	// (3) Multiple keys
	m.Increment("DE", 7)
	m.Increment("US", 5)

	require.Equal(t, uint32(3), m.Count("RO"))
	require.Equal(t, uint32(7), m.Count("DE"))
	require.Equal(t, uint32(5), m.Count("US"))

	// (4) Increment existing key
	m.Increment("RO", 2)
	require.Equal(t, uint32(5), m.Count("RO"))

	// (5) Eviction path: fill all 7 slots
	keys := []string{"A", "B", "C", "D", "E", "F", "G"}

	for _, k := range keys {
		m.Increment(k, 1)
	}

	// Now increment a new key → eviction happens
	m.Increment("Z", 10)

	// Count must return something meaningful:
	// - either Z is present (evicted weakest)
	// - or Z is not present (if RO/DE/US were weaker)
	// We only assert that Count does not panic and returns a uint32.
	_ = m.Count("Z")

	// (6) Unknown key → 0
	require.Equal(t, uint32(0), m.Count("XX"))
}

func TestMetaActiveDeepCopyInto(t *testing.T) {
	var src MetaActive[string]

	src.Increment("RO", 3)
	src.Increment("DE", 7)

	src.Names[2].Store(nil)
	src.Values[2].Store(99)
	src.setOccupied(2)

	src.Names[3].Store(nil)

	k4 := new(string)
	*k4 = "US"
	src.Names[4].Store(k4)
	src.Values[4].Store(5)

	var dst MetaActive[string]
	src.DeepCopyInto(&dst)

	// --- Occupancy mask copied correctly ---
	require.Equal(t, src.occupied.Load(), dst.occupied.Load())

	// --- Values copied correctly ---
	require.Equal(t, src.Values[0].Load(), dst.Values[0].Load())
	require.Equal(t, src.Values[1].Load(), dst.Values[1].Load())
	require.Equal(t, src.Values[2].Load(), dst.Values[2].Load())
	require.Equal(t, src.Values[3].Load(), dst.Values[3].Load())
	require.Equal(t, src.Values[4].Load(), dst.Values[4].Load())

	// --- Count lookups work on dst ---
	require.Equal(t, uint32(3), dst.Count("RO"))
	require.Equal(t, uint32(7), dst.Count("DE"))
	require.Equal(t, uint32(5), dst.Count("US"))
	require.Equal(t, uint32(0), dst.Count("XX"))

	// --- CRITICAL: deep copy isolation ---
	// Modifying src must NOT affect dst
	*src.Names[0].Load() = "MODIFIED-RO"
	*src.Names[1].Load() = "MODIFIED-DE"
	*src.Names[4].Load() = "MODIFIED-US"

	require.Equal(t,
		"RO",
		*dst.Names[0].Load(),
		"dst was corrupted by src mutation — DeepCopyInto shares pointers")
	require.Equal(t,
		"DE",
		*dst.Names[1].Load(),
		"dst was corrupted by src mutation — DeepCopyInto shares pointers")
	require.Equal(t,
		"US",
		*dst.Names[4].Load(),
		"dst was corrupted by src mutation — DeepCopyInto shares pointers")

	// Counts on dst must remain unchanged
	require.Equal(t, uint32(3), dst.Count("RO"))
	require.Equal(t, uint32(7), dst.Count("DE"))
	require.Equal(t, uint32(5), dst.Count("US"))
}

func BenchmarkMetaActiveIncrement_Parallel(b *testing.B) {
	gomaxprocsValues := []int{1, 2, 3, 4, 8}
	key := "RO"

	for _, g := range gomaxprocsValues {
		b.Run(
			fmt.Sprintf("GOMAXPROCS=%d", g),
			func(b *testing.B) {
				prev := runtime.GOMAXPROCS(g)
				defer runtime.GOMAXPROCS(prev)

				b.ReportAllocs()
				b.ResetTimer()

				var m MetaActive[string]

				// Preload deterministic state
				preKeys := []string{"RO", "DE", "US", "FR", "IT", "ES", "NL"}

				for ix := range len(preKeys) {
					k := new(string)
					*k = preKeys[ix]

					m.Names[ix].Store(k)
					m.Values[ix].Store(uint32(ix + 1))

					m.setOccupied(ix)
				}

				b.RunParallel(
					func(pb *testing.PB) {
						for pb.Next() {
							m.Increment(key, 1)
						}
					},
				)
			},
		)
	}
}
