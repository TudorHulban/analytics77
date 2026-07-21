package analytics

import "sync/atomic"

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

type Ingestable[TData any] interface {
	Ingest(TData)
}

type ring[T Ingestable[TData], TData any] struct {
	slots   [7]T
	current atomic.Int32
}

func newRing[T Ingestable[TData], TData any]() *ring[T, TData] {
	var r ring[T, TData]

	r.current.Store(0)

	return &r
}

func NewRegistryRing[T Ingestable[TData], TData any](dstTimestamps ...int64) *RegistryRing[T, TData] {
	result := RegistryRing[T, TData]{
		Ring: newRing[T](),
	}

	if len(dstTimestamps) == 2 {
		result.TimestampDSTSpring = dstTimestamps[0]
		result.TimestampDSTWinter = dstTimestamps[1]
	}

	return &result
}

func (r *ring[T, TData]) ZeroSlot(i int32) {
	var zero T

	r.slots[i] = zero
}

func (r *ring[T, TData]) Advance() {
	next := (r.current.Load() + 1) % int32(len(r.slots))

	r.ZeroSlot(next)

	r.current.Store(next)
}

func (r *ring[T, TData]) GetActiveSlot() *T {
	return &r.slots[r.current.Load()]
}

type RegistryRing[T Ingestable[TData], TData any] struct {
	Ring *ring[T, TData]

	TimestampDSTWinter int64
	TimestampDSTSpring int64

	CalendarMonthCurrentNumber int8 // no need for year as we keep only 7 months.
}

func (rr *RegistryRing[T, TData]) Ingest(data TData) {
	slot := rr.Ring.GetActiveSlot()

	(*slot).Ingest(data)
}
