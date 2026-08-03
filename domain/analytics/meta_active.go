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
	Names  [7]atomic.Pointer[T]
	Values [7]atomic.Uint32

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

				var newVal uint32

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
			m.Values[claimed].Add(byValue)

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

			for {
				oldVal := m.Values[lowestIdx].Load()

				var newVal uint32

				if oldVal > maxUint32-byValue {
					newVal = maxUint32
				} else {
					newVal = oldVal + byValue
				}

				if m.Values[lowestIdx].CompareAndSwap(oldVal, newVal) {
					break
				}
			}

			k := new(T)
			*k = key
			m.Names[lowestIdx].Store(k)

			return
		}
	}
}

func (m *MetaActive[T]) deepCopyInto(dst *MetaActive[T]) {
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

func (m *MetaActive[T]) Count(forKey T) uint32 {
	if m == nil {
		return 0
	}

	mask := m.occupied.Load()

	for ix := range 7 {
		if mask&(1<<ix) == 0 {
			continue // slot not occupied — Names[ix] may be stale, do not trust it
		}

		if ptr := m.Names[ix].Load(); ptr != nil && *ptr == forKey {
			return m.Values[ix].Load()
		}
	}

	return 0
}

func (m *MetaActive[T]) GetValue(byKey T) (uint32, error) {
	mask := m.occupied.Load()

	for ix := range 7 {
		if mask&(1<<ix) == 0 {
			continue // slot not occupied — Names[ix] may be stale, do not trust it
		}

		ptr := m.Names[ix].Load()

		if ptr != nil && *ptr == byKey {
			return m.Values[ix].Load(), nil
		}
	}

	return 0,
		ErrKeyNotFound
}

// MergeFrom accumulates the contents of src into the receiver.
// It sums values for identical keys, then re-sorts and truncates
// to the fixed Top-N capacity (7 entries).
//
// Zero-allocation: both inputs have at most 7 entries each, so the
// merge uses a fixed 14-slot stack array with linear-scan accumulation
// and an insertion sort instead of a map + sort.Slice. For N<=14 this
// is also faster in practice than hashing — no bucket allocation, no
// interface boxing for the comparator closure.
func (m *MetaActive[T]) MergeFrom(src *MetaActive[T]) {
	type kv struct {
		name  T
		value uint32
	}

	var (
		acc   [14]kv
		count int
	)

	accumulate := func(names *[7]atomic.Pointer[T], values *[7]atomic.Uint32) {
		for ix := range 7 {
			namePtr := names[ix].Load()
			if namePtr == nil {
				continue
			}

			val := values[ix].Load()
			name := *namePtr

			found := false

			for j := range count {
				if acc[j].name == name {
					acc[j].value += val
					found = true

					break
				}
			}

			if !found {
				acc[count] = kv{name, val}
				count++
			}
		}
	}

	accumulate(&m.Names, &m.Values)
	accumulate(&src.Names, &src.Values)

	// insertion sort descending by value — trivial cost for count<=14
	for i := 1; i < count; i++ {
		cur := acc[i]
		j := i - 1

		for j >= 0 && acc[j].value < cur.value {
			acc[j+1] = acc[j]
			j--
		}

		acc[j+1] = cur
	}

	// write back top 7
	for ix := range len(m.Names) {
		if ix < count {
			name := acc[ix].name

			m.Names[ix].Store(&name)
			m.Values[ix].Store(acc[ix].value)
		} else {
			m.Names[ix].Store(nil)
			m.Values[ix].Store(0)
		}
	}

	// update occupancy bitmask
	var occ uint32

	for ix := range len(m.Names) {
		if m.Names[ix].Load() != nil {
			occ |= 1 << uint(ix)
		}
	}

	m.occupied.Store(occ)
}

func (m *MetaActive[T]) reset() {
	for i := range m.Names {
		m.Names[i].Store(nil)
		m.Values[i].Store(0)
	}

	m.occupied.Store(0)
}
