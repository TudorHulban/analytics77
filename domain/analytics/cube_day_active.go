package analytics

// OLAP CUBE DESIGN NOTE
//
// This cube is a fixed‑shape, zero‑allocation, cache‑resident analytics kernel.
// It uses 7‑slot Space‑Saving sketches for top‑K tracking and fixed 31×24
// time‑bucket arrays for deterministic traversal.
//
// All structures avoid slices, maps, and heap allocations. Every operation
// (increment, merge, aggregate) runs in strictly bounded time with no GC impact.
//
// IMPORTANT:
// - MetaActive uses a 7‑bit occupancy mask to track valid slots.
//   Zero‑value keys ("" or 0) are valid and must never be treated as empty.
// - MetaArchived is an unordered top‑7 sketch. Index positions are NOT stable.
//   Consumers must use Count(key) instead of assuming ordering.
// - MergeActive and MergeArchived collapse duplicates and select the top 7
//   without sorting, using only fixed arrays and linear scans.
//
// This kernel is designed for deterministic, high‑frequency ingestion and
// predictable OLAP reporting without hidden allocations or dynamic behavior.

type DayActive [24]MetricActive

func (d *DayActive) IsZero() bool {
	for i := range d {
		if !d[i].IsZero() {
			return false
		}
	}

	return true
}
