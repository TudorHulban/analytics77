package analytics

import (
	"fmt"
	"strings"
	"sync/atomic"
)

// OCCUPANCY MASK
//
// MetaActive uses a 7‑bit mask to track which slots contain valid keys.
// Bit i corresponds to slot i. This avoids relying on zero‑value detection,
// which is unsafe because zero‑value keys ("" or 0) are valid.
//
// Example:
//
//	occupied = 0b00101101
//	slots 0, 2, 3, 5 are valid.
//
// This mask is the only correct way to determine slot validity in MetaActive.
type MetaActive[T comparable] struct {
	Names    [7]atomic.Pointer[T]
	Values   [7]atomic.Uint32
	occupied atomic.Uint32 // lower 7 bits used
}

// INCREMENT (MetaActive)
//
// Inserts or updates a key in the 7‑slot Space‑Saving sketch.
// Slot validity is tracked via a 7‑bit occupancy mask. Zero‑value keys
// ("" or 0) are valid and must NOT be treated as empty.
//
// Steps:
// 1. Scan for exact match and increment.
// 2. Scan for an unoccupied slot using the bitmask.
// 3. If full, evict the weakest slot (Space‑Saving) and install the new key.
//
// This function is branch‑predictable, and stable under full load.
// It is the ingestion hot path and must remain minimal and deterministic.
func (m *MetaActive[T]) Increment(key T, byValue uint32) {
outer:
	for {
		mask := m.occupied.Load()

		// 1. Exact-match fast path — CAS-based, no Add
		for ix := range 7 {
			if (mask & (1 << ix)) == 0 {
				continue
			}

			ptr := m.Names[ix].Load()
			if ptr == nil || *ptr != key {
				continue
			}

			// Key matched
			// Now CAS the value while verifying key has not changed
			for {
				oldVal := m.Values[ix].Load()

				// Re-verify key before every CAS attempt
				if after := m.Names[ix].Load(); after == nil || *after != key {
					continue outer // key changed, retry from scratch
				}

				newVal := oldVal //nolint:ineffassign
				if oldVal > maxUint32-byValue {
					newVal = maxUint32
				} else {
					newVal = oldVal + byValue
				}

				if m.Values[ix].CompareAndSwap(oldVal, newVal) {
					return // success
				}
				// CAS failed, loop back and re-verify key + retry
			}
		}

		// 2. Claim an empty slot — unchanged
		for {
			claimed := -1

			for ix := range 7 {
				if (mask & (1 << ix)) != 0 {
					continue
				}

				newMask := mask | (1 << ix)
				if m.occupied.CompareAndSwap(mask, newMask) {
					claimed = ix

					break
				}

				mask = m.occupied.Load()
			}

			if claimed == -1 {
				if mask&0b1111111 == 0b1111111 {
					break
				}

				continue
			}

			k := new(T)
			*k = key
			m.Names[claimed].Store(k)
			m.Values[claimed].Store(byValue)

			return
		}

		// 3. Full: evict the weakest slot — unchanged
		for {
			lowestIdx := 0
			lowestVal := m.Values[0].Load()

			for ix := 1; ix < 7; ix++ {
				if v := m.Values[ix].Load(); v < lowestVal {
					lowestVal = v
					lowestIdx = ix
				}
			}

			newVal := lowestVal //nolint:ineffassign
			if lowestVal > maxUint32-byValue {
				newVal = maxUint32
			} else {
				newVal = lowestVal + byValue
			}

			if !m.Values[lowestIdx].CompareAndSwap(lowestVal, newVal) {
				continue
			}

			k := new(T)
			*k = key
			m.Names[lowestIdx].Store(k)

			return
		}
	}
}

// AsMetaArchive converts a MetaActive[T] into a MetaArchived[T].
// It performs a deterministic top‑7 extraction with zero allocations.
//
// WARNING: This function assumes no concurrent writes to m.
func (m *MetaActive[T]) AsMetaArchive() MetaArchived[T] {
	var (
		tmpKeys [7]T
		tmpVals [7]uint32
		tmpLen  int
	)

	// 1. Load active entries into temporary arrays
	mask := m.occupied.Load()

	for ix := range 7 {
		if (mask & (1 << ix)) != 0 {
			if ptr := m.Names[ix].Load(); ptr != nil {
				tmpKeys[tmpLen] = *ptr
				tmpVals[tmpLen] = m.Values[ix].Load()

				tmpLen++
			}
		}
	}

	// 2. Combine duplicates (Space‑Saving can produce same key in multiple slots)
	for ix := 0; ix < tmpLen; ix++ {
		if tmpVals[ix] == 0 {
			continue
		}

		ki := tmpKeys[ix]

		for j := ix + 1; j < tmpLen; j++ {
			if tmpVals[j] != 0 && tmpKeys[j] == ki {
				tmpVals[ix] = tmpVals[ix] + tmpVals[j]
				tmpVals[j] = 0
			}
		}
	}

	// 3. Select top 7 by value (no sorting, zero alloc)
	var out MetaArchived[T]

	for k := range 7 {
		bestIdx := -1
		bestVal := uint32(0)

		for i := 0; i < tmpLen; i++ {
			v := tmpVals[i]
			if v > bestVal {
				bestVal = v
				bestIdx = i
			}
		}

		if bestIdx == -1 || bestVal == 0 {
			break
		}

		out.Names[k] = tmpKeys[bestIdx] //nolint:gosec
		out.Values[k] = bestVal

		tmpVals[bestIdx] = 0 //nolint:gosec
	}

	return out
}

func (m *MetaActive[T]) DeepCopyInto(dst *MetaActive[T]) {
	// 1. Copy non-atomic fields
	dst.occupied.Store(m.occupied.Load())

	// 2. Copy atomic pointers safely
	for ix := range 7 {
		ptr := m.Names[ix].Load()

		if ptr != nil {
			k := new(T)
			*k = *ptr

			dst.Names[ix].Store(k)
		}
	}

	// 3. Copy atomic counters safely
	for ix := range 7 {
		val := m.Values[ix].Load()
		dst.Values[ix].Store(val)
	}
}

// If occupied is 0, we trust that Names and Values are also empty.
func (m *MetaActive[T]) IsZero() bool {
	return m.occupied.Load() == 0
}

func (m *MetaActive[T]) String() string {
	var b strings.Builder

	mask := m.occupied.Load()

	for ix := range 7 {
		if (mask & (1 << ix)) == 0 {
			continue
		}

		if ptr := m.Names[ix].Load(); ptr != nil {
			fmt.Fprintf(
				&b,
				"%v:%d\n",

				*ptr,
				m.Values[ix].Load(),
			)
		}
	}

	return b.String()
}

func (m *MetaActive[T]) Count(key T) uint32 {
	if m == nil {
		return 0
	}

	for ix := range 7 {
		if ptr := m.Names[ix].Load(); ptr != nil && *ptr == key {
			return m.Values[ix].Load()
		}
	}

	return 0
}

func (m *MetaActive[T]) GetValue(byKey T) (uint32, error) {
	for ix := range 7 {
		ptr := m.Names[ix].Load()

		if ptr != nil && *ptr == byKey {
			return m.Values[ix].Load(), nil
		}
	}

	return 0, ErrKeyNotFound
}
