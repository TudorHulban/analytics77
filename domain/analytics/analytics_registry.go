package analytics

import "sync/atomic"

// Registry keeps info as per UTC times.
type Registry struct {
	Slots       [7]MonthActive
	CurrentSlot atomic.Int32

	TimestampDSTWinter int64
	TimestampDSTSpring int64

	CalendarMonthCurrentNumber int8 // no need for year as we keep only 7 months.
}

func NewRegistry(dstTimestamps ...int64) *Registry {
	result := Registry{}

	if len(dstTimestamps) == 2 {
		result.TimestampDSTSpring = dstTimestamps[0]
		result.TimestampDSTWinter = dstTimestamps[1]
	}

	return &result
}

func (r *Registry) zeroSlot(slotNo int32) {
	var zero MonthActive

	r.Slots[slotNo] = zero
}

func (r *Registry) Advance() {
	next := (r.CurrentSlot.Load() + 1) % int32(len(r.Slots))

	r.zeroSlot(next)

	r.CurrentSlot.Store(next)
}

// CurrentMonth returns the active buffer for the current month.
// Callers must not cache the returned pointer beyond a single logical
// operation — always re-call CurrentMonth() for the next one.
func (r *Registry) GetActiveSlot() *MonthActive {
	return &r.Slots[r.CurrentSlot.Load()]
}

// PreviousMonth returns the active buffer for the previous month.
// Same caching rule as CurrentMonth.
func (r *Registry) GetPreviousSlot() *MonthActive {
	prev := (r.CurrentSlot.Load() + int32(len(r.Slots)) - 1) % int32(len(r.Slots))

	return &r.Slots[prev]
}

func (r *Registry) ForEachMetric(slot *MonthActive, fn func(day int8, hour int8, m *MetricActive)) {
	for day := int8(0); day < 31; day++ {
		for hour := int8(0); hour < 24; hour++ {
			metric := (*slot).GetMetric(day, hour)

			fn(day, hour, metric)
		}
	}
}
