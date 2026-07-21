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

type TMetric interface {
	GetRecordsPerPeriod() uint32
}

type Ingestable[Metric TMetric] interface {
	Ingest(Metric)
	GetMetric(day int8, hour int8) *Metric
}

type ring[T Ingestable[Data], Data TMetric] struct {
	slots   [7]T
	current atomic.Int32
}

func newRing[T Ingestable[Data], Data TMetric]() *ring[T, Data] {
	var r ring[T, Data]

	r.current.Store(0)

	return &r
}

func (r *ring[T, Data]) ZeroSlot(i int32) {
	var zero T

	r.slots[i] = zero
}

func (r *ring[T, Data]) Advance() {
	next := (r.current.Load() + 1) % int32(len(r.slots))

	r.ZeroSlot(next)

	r.current.Store(next)
}

func (r *ring[T, Data]) GetActiveSlot() *T {
	return &r.slots[r.current.Load()]
}

func (r *ring[T, Data]) GetPreviousSlot() *T {
	curr := r.current.Load()
	prev := (curr + int32(len(r.slots)) - 1) % int32(len(r.slots))

	return &r.slots[prev]
}

type Registry[T Ingestable[Data], Data TMetric] struct {
	Ring *ring[T, Data]

	TimestampDSTWinter int64
	TimestampDSTSpring int64

	CalendarMonthCurrentNumber int8 // no need for year as we keep only 7 months.
}

func NewRegistry[T Ingestable[Data], Data TMetric](dstTimestamps ...int64) *Registry[T, Data] {
	result := Registry[T, Data]{
		Ring: newRing[T](),
	}

	if len(dstTimestamps) == 2 {
		result.TimestampDSTSpring = dstTimestamps[0]
		result.TimestampDSTWinter = dstTimestamps[1]
	}

	return &result
}

func (r *Registry[T, Data]) Ingest(data Data) {
	slot := r.Ring.GetActiveSlot()

	(*slot).Ingest(data)
}

func (r *Registry[T, Data]) ForEachMetric(slot *T, fn func(day int8, hour int8, m *Data)) {
	for day := int8(0); day < 31; day++ {
		for hour := int8(0); hour < 24; hour++ {
			m := (*slot).GetMetric(day, hour)

			fn(day, hour, m)
		}
	}
}
