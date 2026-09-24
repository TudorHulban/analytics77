package analytics

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)

// Registry keeps info as per UTC times.
type Registry struct {
	Slots [7]*MonthActive

	TimestampDSTWinter int64
	TimestampDSTSpring int64

	CurrentSlot                atomic.Int32
	CalendarMonthCurrentNumber atomic.Int32

	muAdvance sync.Mutex
}

// NewRegistry initial slots:
//
// Slots[0] = current month
//
// Slots[1] = 1 month back
//
// Slots[2] = 2 months back
//
// Slots[3] = 3 months back
//
// Slots[4] = 4 months back
//
// Slots[5] = 5 months back
//
// Slots[6] = 6 months back
func NewRegistry(inMonth int8, dstTimestamps ...int64) *Registry {
	var month int32

	if inMonth < 1 {
		month = 1
	} else if inMonth > 12 {
		month = 12
	} else {
		month = int32(inMonth)
	}

	result := Registry{}
	result.CalendarMonthCurrentNumber.Store(month)

	for i := range result.Slots {
		result.Slots[i] = &MonthActive{}
	}

	if len(dstTimestamps) == 2 {
		result.TimestampDSTSpring = dstTimestamps[0]
		result.TimestampDSTWinter = dstTimestamps[1]
	}

	return &result
}

func (r *Registry) zeroSlot(slotNo int32) {
	slot := r.Slots[slotNo]

	for day := range int8(31) {
		for hour := range int8(24) {
			m := &slot[day][hour]
			m.RecordsPerPeriod.Store(0)

			m.TopIPs.reset()
			m.TopASN.reset()
			m.TopCountries.reset()
			m.TopCities.reset()
			m.TopURLs.reset()
			m.TopOperatingSystems.reset()
			m.TopBrowsers.reset()
		}
	}
}

// CurrentMonth returns the active buffer for the current month.
// Callers must not cache the returned pointer beyond a single logical
// operation — always re-call CurrentMonth() for the next one.
func (r *Registry) GetActiveSlot() *MonthActive {
	return r.Slots[r.CurrentSlot.Load()]
}

// PreviousMonth returns the active buffer for the previous month.
// Same caching rule as CurrentMonth.
func (r *Registry) GetPreviousSlot() *MonthActive {
	prev := (r.CurrentSlot.Load() + int32(len(r.Slots)) - 1) % int32(len(r.Slots))

	return r.Slots[prev]
}

func (*Registry) ForEachMetric(slot *MonthActive, fn func(day int8, hour int8, m *MetricActive)) {
	for day := range int8(31) {
		for hour := range int8(24) {
			metric := (*slot).getMetric(day, hour)

			fn(day, hour, metric)
		}
	}
}

func (r *Registry) WriteTo(b *strings.Builder) {
	r.GetActiveSlot().WriteTo(b, "current")
	r.GetPreviousSlot().WriteTo(b, "previous")

	// archived months (the remaining 5 slots)
	for ix := range 7 {
		month := r.Slots[ix]

		// skip current and previous
		if ix == int(r.CurrentSlot.Load()) {
			continue
		}

		if ix == int((r.CurrentSlot.Load()+7-1)%7) {
			continue
		}

		month.WriteTo(b, fmt.Sprintf("slot[%d]", ix))
	}
}
