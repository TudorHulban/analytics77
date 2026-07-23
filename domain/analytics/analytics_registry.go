package analytics

import "sync/atomic"

// Registry keeps info as per UTC times.
type Registry struct {
	Slots [7]MonthActive

	TimestampDSTWinter int64
	TimestampDSTSpring int64

	CurrentSlot                atomic.Int32
	CalendarMonthCurrentNumber int8
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
	slot := &r.Slots[slotNo]

	var (
		defaultString  string
		defaultOS      OS
		defaultBrowser Browser
	)

	for day := range 31 {
		for hour := range 24 {
			m := &slot[day][hour]

			// reset atomic counter
			m.RecordsPerPeriod.Store(0)

			// reset TopIPs
			for i := range m.TopIPs.Names {
				m.TopIPs.Names[i].Store(&defaultString)
				m.TopIPs.Values[i].Store(0)
			}

			// reset TopCountries
			for i := range m.TopCountries.Names {
				m.TopCountries.Names[i].Store(&defaultString)
				m.TopCountries.Values[i].Store(0)
			}

			// reset TopASN
			for i := range m.TopASN.Names {
				m.TopASN.Names[i].Store(&defaultString)
				m.TopASN.Values[i].Store(0)
			}

			// reset TopCities
			for i := range m.TopCities.Names {
				m.TopCities.Names[i].Store(&defaultString)
				m.TopCities.Values[i].Store(0)
			}

			// reset TopURL
			for i := range m.TopURLs.Names {
				m.TopURLs.Names[i].Store(&defaultString)
				m.TopURLs.Values[i].Store(0)
			}

			// reset TopOperatingSystems
			for i := range m.TopOperatingSystems.Names {
				m.TopOperatingSystems.Names[i].Store(&defaultOS)
				m.TopOperatingSystems.Values[i].Store(0)
			}

			// reset TopBrowsers
			for i := range m.TopBrowsers.Names {
				m.TopBrowsers.Names[i].Store(&defaultBrowser)
				m.TopBrowsers.Values[i].Store(0)
			}
		}
	}
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

func (*Registry) ForEachMetric(slot *MonthActive, fn func(day int8, hour int8, m *MetricActive)) {
	for day := range int8(31) {
		for hour := range int8(24) {
			metric := (*slot).GetMetric(day, hour)

			fn(day, hour, metric)
		}
	}
}
