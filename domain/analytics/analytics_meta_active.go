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

		// 1. Exact-match fast path.
		// A slightly stale mask is fine here: bits only ever transition 0->1
		// outside of Rollover/DeepCopy (which never run concurrently with
		// ingestion), so we can't see a "phantom" occupied slot that isn't
		// really there.
		for ix := range 7 {
			if (mask & (1 << ix)) == 0 {
				continue
			}

			ptr := m.Names[ix].Load()
			if ptr == nil || *ptr != key {
				continue
			}

			m.Values[ix].Add(byValue)

			// Guard: if a concurrent eviction replaced this slot's key while
			// above Add was in flight, the increment just landed on the wrong
			// bucket. Undo it and retry the whole operation from scratch.
			if after := m.Names[ix].Load(); after == nil || *after != key {
				m.Values[ix].Add(-byValue) // unsigned wraparound == subtract

				continue outer
			}

			return
		}

		// 2. Claim an empty slot.
		// Each bit is claimed with its own CAS, so only one goroutine can ever
		// win a given slot. A failed CAS means someone else just claimed a
		// bit (this one or another); reload and keep scanning forward rather
		// than retrying the same bit.
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
					break // all 7 slots occupied -> fall through to eviction
				}

				continue // fresh mask picked up a newly-freed view; rescan
			}

			// We now exclusively own `claimed`: occupancy bits only ever go
			// 0->1 via this CAS, so no other goroutine can also hold it.
			k := new(T)
			*k = key

			m.Names[claimed].Store(k)
			m.Values[claimed].Store(byValue)

			return
		}

		// 3. Full: evict the weakest slot (Space-Saving).
		// The CAS on Values makes "read weakest, then update it" atomic
		// against other evictions and other exact-match Adds racing the same
		// slot: if anything touches Values[lowestIdx] between our read and
		// our swap, the CAS fails and we just recompute the weakest slot.
		for {
			lowestIdx := 0
			lowestVal := m.Values[0].Load()

			for ix := 1; ix < 7; ix++ {
				if v := m.Values[ix].Load(); v < lowestVal {
					lowestVal = v
					lowestIdx = ix
				}
			}

			newVal := lowestVal + byValue

			if !m.Values[lowestIdx].CompareAndSwap(lowestVal, newVal) {
				continue // lost the race on this slot's value; retry
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
	for ix := range 7 {
		if (m.occupied.Load() & (1 << ix)) != 0 {
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
		dst.Names[ix].Store(ptr)
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

func (m *MetaActive[T]) setOccupied(ix int) {
	bit := uint32(1 << ix)

	for {
		oldValue := m.occupied.Load()
		newValue := oldValue | bit

		if m.occupied.CompareAndSwap(oldValue, newValue) {
			return
		}
	}
}
