package analytics

/*
Package analytics implements a zero‑allocation OLAP cube for high‑throughput
telemetry ingestion and deterministic, bounded‑time queries.

This is systems‑level code: fixed‑shape arrays, atomic increments, no maps,
no slices in hot paths, no GC pressure, and no interface indirection.

-------------------------------------------------------------------------------
Core Concepts
-------------------------------------------------------------------------------

Registry
    Represents all analytics for a single Site.
    Contains:
        - MonthCurrent:   31 days × 24 hours of active metrics
        - MonthPrevious:  previous month’s active metrics
        - History:        7 months of archived, merged metrics

    All structures are fixed‑size arrays. No dynamic allocation occurs during
    ingestion or queries.

MetricActive
    The ingestion‑hot structure:
        - atomic counter (RecordsPerPeriod)
        - seven MetaActive top‑K sketches
        - zero allocations
        - safe for lock‑free increments

MetricArchived
    The compact, merged representation used for history retention.

-------------------------------------------------------------------------------
Ingestion Model
-------------------------------------------------------------------------------

Ingestion is performed through the DataCenter:
    - structural mutations (site creation)
    - metric increments are lock‑free and use atomics
    - no allocations occur during ingestion
    - day indexing uses explicit helpers:
          CalendarDayToIndex(1–31 → 0–30)
          IndexToCalendarDay(0–30 → 1–31)

This ensures predictable performance under load.

-------------------------------------------------------------------------------
Query Model
-------------------------------------------------------------------------------

Query methods return fixed arrays + count, never slices:
    days, n := r.CurrentMonthDaysWithData()
    hours, m := r.CurrentMonthHoursWithData(calendar day)

This guarantees:
    - zero allocations
    - no slice escapes
    - deterministic memory behavior

Presence scans are bounded:
    - max 31 days
    - max 24 hours per day

-------------------------------------------------------------------------------
History Retention
-------------------------------------------------------------------------------

The cube maintains:
    - current month (active)
    - previous month (active)
    - 7 archived months (merged)

Month rollover performs:
    1. merge active → archived
    2. clear active
    3. shift history window

All operations are deterministic and bounded.

-------------------------------------------------------------------------------
Day Index Helpers
-------------------------------------------------------------------------------

Internal indexing uses 0–30.
Calendar days use 1–31.

Two pure helpers define the contract:
    CalendarDayToIndex(day) → day‑1
    IndexToCalendarDay(ix)  → ix+1

These eliminate off‑by‑one errors across ingestion, queries, and history.

-------------------------------------------------------------------------------
Testing Philosophy
-------------------------------------------------------------------------------

Tests operate on:
    - fixed arrays
    - returned counts
    - deterministic presence scans

Example:
    hours, n := r.CurrentMonthHoursWithData(4)
    require.Equal(t, n, 3)
    require.ElementsMatch(t, hours[:n], []int8{1, 7, 22})

Tests do not rely on slice lengths, because arrays are fixed‑size.

-------------------------------------------------------------------------------
Design Goals
-------------------------------------------------------------------------------

    - Zero allocations in ingestion and queries
    - Lock‑free increments
    - Deterministic memory layout
    - Bounded time operations
    - No GC pressure
    - No interfaces, no mocks
      (this is systems code, not application code)

This package is intended for high‑throughput telemetry pipelines, embedded
analytics, and real‑time observability systems.

-------------------------------------------------------------------------------
When to Use This Package
-------------------------------------------------------------------------------

Use it when you need:
    - millions of events per minute
    - predictable latency
    - fixed memory footprint
    - OLAP‑style rollups
    - top‑K per hour/day/month
    - zero GC interference

Do not use it for:
    - ad‑hoc analytics
    - arbitrary dimensions
    - dynamic schemas
    - general‑purpose BI workloads
*/

// Ring invariants:
//
// 1. The ring contains N fixed slots. Slots are never added, removed, or reordered.
//    Their memory addresses remain stable for the lifetime of the Ring.
//
// 2. Exactly one slot is "active" at any time. The active slot index is stored in
//    r.current and read atomically.
//
// 3. Ingestion writes ONLY to the active slot. No ingestion touches any other slot.
//
// 4. All non-active slots are treated as read-only history. Queries may read any
//    history slot, but MUST NOT read the active slot unless explicitly performing
//    a "current" query.
//
// 5. Switching the active slot is done by atomically updating r.current. No other
//    synchronization is required. Slot switching never races with ingestion because
//    ingestion always reads r.current before writing.
//
// 6. A slot becomes the new active slot only after it has been fully zeroed or
//    reset. Zeroing MUST occur before the atomic index update.
//
// 7. No method may mix reads from the active slot and history slots. Query semantics
//    must be strictly either "current" or "history", never both.
//
// 8. Slots store values (not pointers). This guarantees stable memory layout,
//    eliminates pointer races, improves cache locality, and prevents accidental
//    slot replacement.
//
// These invariants make the Ring a lock-free, time-partitioned structure with
// deterministic ingestion and query behavior.
