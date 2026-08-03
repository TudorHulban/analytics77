package analytics

// Package analytics implements a zero-allocation OLAP cube for high-throughput
// telemetry ingestion and deterministic, bounded-time queries.
//
// This is systems-level code: fixed-shape arrays, atomic increments, no maps,
// no slices in hot paths, no GC pressure, and no interface indirection.
//
// # Core Concepts
//
// Registry
// The main entry point representing all analytics for a single site.
// It maintains a 7-slot ring buffer of *MonthActive structures:
//   - Slots[CurrentSlot]: The active, writable month.
//   - Slots[1..6]: Read-only history months (previous month, 2 months back, etc.).
//
// All structures are fixed-size arrays. No dynamic allocation occurs during
// ingestion or queries.
//
// MetricActive
// The ingestion-hot structure representing a single hour's data. It contains:
//   - RecordsPerPeriod: An atomic counter for total records.
//   - Seven MetaActive top-K sketches: TopIPs, TopASN, TopCountries, TopCities,
//     TopURLs, TopOperatingSystems, and TopBrowsers.
//
// MetaActive[T comparable]
// A 7-slot Space-Saving sketch for top-K tracking. It uses a 7-bit occupancy
// mask to track valid slots, ensuring zero-value keys ("" or 0) are treated
// as valid and never mistaken for empty slots. It provides lock-free Increment,
// DeepCopyInto, and MergeFrom operations.
//
// # Ingestion Model
//
// Ingestion is designed for maximum throughput:
//   - Metric increments are lock-free and use atomic operations.
//   - No heap allocations occur during ingestion.
//   - Day indexing uses explicit 0–30 bounds.
//   - Writes target ONLY the active slot (Slots[CurrentSlot]).
//
// # Query Model
//
// Query methods return fixed arrays or aggregated structures, never dynamically
// allocated slices, guaranteeing:
//   - Zero allocations.
//   - No slice escapes.
//   - Deterministic memory behavior.
//
// Aggregation is available at multiple granularities:
//   - Hourly: CurrentMonthAggregateTopNForHour, PreviousMonthAggregateTopNForHour
//   - Daily: CurrentMonthAggregateTopNForDay, PreviousMonthAggregateTopNForDay
//   - Monthly: CurrentMonthAggregateTopN, PreviousMonthAggregateTopN
//   - Historical: HistoryAggregateTopN, HistoryAggregateTopNForMonth, etc.
//
// # Time Management & Rollover
//
// The Registry uses a ring buffer for month rollover. Calling Advance():
//  1. Calculates the next slot index: (CurrentSlot + 1) % 7.
//  2. Zeroes the target slot completely (resetting all atomic counters and masks).
//  3. Atomically updates CurrentSlot to the new index.
//  4. Increments CalendarMonthCurrentNumber (wrapping from 12 to 1 via CAS).
//
// This ensures the new active slot is clean before it becomes writable,
// and the switch is atomic and lock-free.
//
// # Supporting Types
//
//   - GeoIP: Represents parsed geolocation and ASN data for an IP address.
//   - AggregatedTopN: A container holding merged MetaActive sketches for
//     query responses, combining data across multiple time buckets.
//   - OS / Browser: Comparable key types used in their respective Top-K sketches.
//
// # Design Goals
//
//   - Zero allocations in ingestion and queries.
//   - Lock-free increments via atomic operations.
//   - Deterministic memory layout and cache locality.
//   - Bounded-time operations (strictly O(1) or small fixed N).
//   - No GC pressure.
//   - No interfaces, no mocks (this is systems code, not application code).
//
// # When to Use This Package
//
// Use it when you need:
//   - Millions of events per minute.
//   - Predictable, bounded latency.
//   - Fixed memory footprint.
//   - OLAP-style rollups and top-K per hour/day/month.
//   - Zero GC interference.
//
// Do not use it for:
//   - Ad-hoc analytics with arbitrary dimensions.
//   - Dynamic schemas.
//   - General-purpose BI workloads requiring flexible querying.
//
// # Ring Invariants
//
//  1. The ring contains exactly 7 fixed slots. Slots are never added, removed,
//     or reordered. Their memory addresses remain stable for the lifetime of the Registry.
//  2. Exactly one slot is "active" at any time. The active slot index is stored in
//     CurrentSlot and read atomically.
//  3. Ingestion writes ONLY to the active slot. No ingestion touches any other slot.
//  4. All non-active slots are treated as read-only history. Queries may read any
//     history slot, but MUST NOT read the active slot unless explicitly performing
//     a "current" query.
//  5. Switching the active slot is done by atomically updating CurrentSlot. No other
//     synchronization is required. Slot switching never races with ingestion because
//     ingestion always reads CurrentSlot before writing.
//  6. A slot becomes the new active slot only after it has been fully zeroed.
//     Zeroing MUST occur before the atomic index update.
//  7. No method may mix reads from the active slot and history slots in a single
//     logical operation. Query semantics must be strictly either "current" or "history".
//  8. Slots store atomic pointers to values. This guarantees stable memory layout,
//     eliminates pointer races, improves cache locality, and prevents accidental
//     slot replacement.
