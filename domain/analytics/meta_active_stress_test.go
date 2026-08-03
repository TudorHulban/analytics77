package analytics

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestMetaActiveIncrement_ConcurrentHighCardinality forces every code path in
// Increment under real concurrency: exact-match, empty-slot claim, AND
// eviction — unlike BenchmarkMetaActiveIncrement_Parallel, which only ever
// touches one pre-warmed key and never exercises claim/eviction races.
//
// Invariant under test: every Increment(_, 1) call adds exactly 1 to the sum
// of all Values, regardless of which branch it takes:
//   - exact match:      Values[ix] += 1
//   - empty slot claim:  a new slot is created holding 1
//   - eviction:          Values[lowestIdx] becomes lowestVal + 1
//
// So after N total calls (byValue=1 each), sum(Values) must equal exactly N.
// A lost update (two goroutines claiming/evicting the same slot) would make
// the sum come out LOWER than N.
func TestMetaActiveIncrement_ConcurrentHighCardinality(t *testing.T) {
	const (
		numGoroutines     = 64
		incrementsPerGoro = 2000
		numDistinctKeys   = 37 // >> 7 slots, guarantees constant eviction churn
	)

	var m MetaActive[string]

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for g := range numGoroutines {
		go func(seed int) {
			defer wg.Done()

			for ix := range incrementsPerGoro {
				key := fmt.Sprintf(
					"k%d",
					(seed*incrementsPerGoro+ix)%numDistinctKeys,
				)

				m.Increment(key, 1)
			}
		}(g)
	}

	wg.Wait()

	var sum uint64

	for ix := range 7 {
		sum = sum + uint64(m.Values[ix].Load())
	}

	require.EqualValues(t,
		numGoroutines*incrementsPerGoro,
		sum,

		"sum of Values must equal total increment calls; a lower sum (%d vs %d) means a lost update (race) in Increment",
		sum,
		numGoroutines*incrementsPerGoro,
	)

	// Sanity: occupancy mask must show all 7 slots full given >>7 distinct keys.
	require.EqualValues(t,
		uint8(0b1111111),
		m.occupied.Load()&0b1111111,
	)
}

// TestMetaActiveIncrement_ConcurrentSameKeys stresses the exact-match path
// itself under contention from many goroutines hammering a small, fixed set
// of keys (no evictions), verifying no lost updates there either.
func TestMetaActiveIncrement_ConcurrentSameKeys(t *testing.T) {
	const (
		numGoroutines     = 64
		incrementsPerGoro = 2000
	)

	keys := []string{"RO", "DE", "US", "FR", "IT", "ES", "NL"}

	var m MetaActive[string]

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for g := range numGoroutines {
		go func(seed int) {
			defer wg.Done()

			for ix := range incrementsPerGoro {
				m.Increment(keys[(seed+ix)%len(keys)], 1)
			}
		}(g)
	}

	wg.Wait()

	var sum uint64

	for ix := range 7 {
		sum = sum + uint64(m.Values[ix].Load())
	}

	require.EqualValues(t,
		numGoroutines*incrementsPerGoro,
		sum,

		"expected: %d, got: %d",
		numGoroutines*incrementsPerGoro,
		sum,
	)
}
