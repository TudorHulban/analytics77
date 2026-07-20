package analytics

import (
	"fmt"
	"strings"
)

type MetaArchived[T comparable] struct {
	Names  [7]T
	Values [7]uint32
}

func (m MetaArchived[T]) IsZero() bool {
	for _, value := range m.Values {
		if value != 0 {
			return false
		}
	}

	return true
}

func (m MetaArchived[T]) String() string {
	var b strings.Builder

	for i := range 7 {
		name := m.Names[i]
		val := m.Values[i]

		if val == 0 {
			continue
		}

		fmt.Fprintf(&b, "%v:%d\n", name, val)
	}

	return b.String()
}

// MERGE ACTIVE (MetaActive → MetaArchived)
//
// Merges a 7‑slot active sketch into an archived 7‑slot sketch.
// Both sketches may contain zero‑value keys ("" or 0), so we rely solely on the
// occupancy mask (MetaActive) and non‑zero values (MetaArchived) to determine
// slot validity. Zero‑value key detection MUST NOT be used.
//
// Steps:
//  1. Load dst (archived) entries into a 14‑slot temp buffer.
//  2. Load src (active) entries using the occupancy bitmask.
//  3. Collapse duplicates by summing values; consumed entries are marked by
//     setting tmpVals[j] = 0. tmpKeys[j] is ignored after that.
//  4. Select the top 7 entries by repeated linear max‑scan (no sorting).
//
// This function is fully deterministic, zero‑alloc, and stable under full load.
// Ordering is NOT guaranteed; callers must use Count(key).
func (m MetaArchived[T]) MergeActive(src *MetaActive[T]) MetaArchived[T] {
	// Temporary fixed array for up to 14 entries
	var (
		tmpKeys [14]T
		tmpVals [14]uint32
		tmpLen  int
	)

	// 1. Load archived entries
	for ix := range 7 {
		value := m.Values[ix]

		if value != 0 {
			tmpKeys[tmpLen] = m.Names[ix]
			tmpVals[tmpLen] = value

			tmpLen++
		}
	}

	// 2. Load active entries
	for i := range 7 {
		if (src.occupied.Load() & (1 << i)) != 0 {
			if ptr := src.Names[i].Load(); ptr != nil {
				tmpKeys[tmpLen] = *ptr
				tmpVals[tmpLen] = src.Values[i].Load()

				tmpLen++
			}
		}
	}

	// 3. Combine duplicates (still zero alloc)
	for i := 0; i < tmpLen; i++ {
		if tmpVals[i] == 0 {
			continue
		}

		ki := tmpKeys[i]

		for j := i + 1; j < tmpLen; j++ {
			if tmpVals[j] != 0 && tmpKeys[j] == ki {
				if tmpVals[i] > maxUint32-tmpVals[j] {
					tmpVals[i] = maxUint32
				} else {
					tmpVals[i] = tmpVals[i] + tmpVals[j]
				}

				tmpVals[j] = 0 // mark as consumed
			}
		}
	}

	// 4. Select top 7 (no sorting, no allocs)
	var out MetaArchived[T]

	for k := range 7 {
		bestIdx := -1
		bestVal := uint32(0)

		for ix := 0; ix < tmpLen; ix++ {
			v := tmpVals[ix]

			if v > bestVal {
				bestVal = v
				bestIdx = ix
			}
		}

		if bestIdx == -1 {
			break
		}

		out.Names[k] = tmpKeys[bestIdx] //nolint:gosec
		out.Values[k] = bestVal

		tmpVals[bestIdx] = 0 //nolint:gosec
	}

	return out
}

// MERGE ARCHIVED (MetaArchived → MetaArchived)
//
// Merges two archived 7‑slot sketches into a new top‑7 result.
// Archived sketches do not track occupancy; slot validity is determined solely
// by non‑zero Values[ix]. Zero‑value keys ("" or 0) are valid and must not be
// treated as empty.
//
// Steps:
//  1. Load all non‑zero entries from dst and src into a 14‑slot temp buffer.
//  2. Collapse duplicates by summing values; consumed entries are marked by
//     setting tmpVals[j] = 0. tmpKeys[j] is ignored after that.
//  3. Select the top 7 entries via repeated linear max‑scan.
//
// This is a pure top‑K merge: no ordering guarantees, no allocations, and
// strictly bounded runtime. Consumers must use Count(key) to query results.
func (m MetaArchived[T]) MergeArchived(src MetaArchived[T]) MetaArchived[T] {
	// Temporary fixed array for up to 14 entries
	var (
		tmpKeys [14]T
		tmpVals [14]uint32

		tmpLen int
	)

	// Load archived entries from dst
	for ix := range 7 {
		value := m.Values[ix]

		if value != 0 {
			tmpKeys[tmpLen] = m.Names[ix]
			tmpVals[tmpLen] = value

			tmpLen++
		}
	}

	// Load archived entries from src
	for ix := range 7 {
		value := src.Values[ix]

		if value != 0 {
			tmpKeys[tmpLen] = src.Names[ix]
			tmpVals[tmpLen] = value

			tmpLen++
		}
	}

	// Combine duplicates (still zero alloc)
	for i := 0; i < tmpLen; i++ {
		if tmpVals[i] == 0 { // skip consumed entries
			continue
		}

		ki := tmpKeys[i]

		for j := i + 1; j < tmpLen; j++ {
			if tmpVals[j] != 0 && tmpKeys[j] == ki {
				if tmpVals[i] > maxUint32-tmpVals[j] {
					tmpVals[i] = maxUint32
				} else {
					tmpVals[i] = tmpVals[i] + tmpVals[j]
				}

				tmpVals[j] = 0 // mark as consumed
			}
		}
	}

	// 4. Select top 7 (no sorting, no allocs)
	var out MetaArchived[T]

	for k := range 7 {
		bestIdx := -1
		bestVal := uint32(0)

		for ix := 0; ix < tmpLen; ix++ {
			value := tmpVals[ix]

			if value > bestVal {
				bestVal = value
				bestIdx = ix
			}
		}

		if bestIdx == -1 {
			break
		}

		out.Names[k] = tmpKeys[bestIdx] //nolint:gosec
		out.Values[k] = bestVal

		tmpVals[bestIdx] = 0 //nolint:gosec
	}

	return out
}

// COUNT LOOKUP
//
// MetaArchived is an unordered top‑7 sketch. Slot positions are not stable.
// Count(key) performs a linear scan over the 7 slots and returns the value
// associated with the key, or 0 if not present.
//
// This is the only correct way to query archived results.
func (m *MetaArchived[T]) Count(forName T) uint32 {
	if m == nil {
		return 0
	}

	for ix := range 7 {
		if m.Names[ix] == forName {
			return m.Values[ix]
		}
	}

	return 0
}
