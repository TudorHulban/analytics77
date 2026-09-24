package analytics

import (
	"fmt"
	"strings"
)

type MonthActive [31]DayActive

func (m *MonthActive) IsZero() bool {
	for i := range m {
		if !m[i].IsZero() {
			return false // Found data in one of the days
		}
	}

	return true // All 31 days are completely zero
}

func (m *MonthActive) deepCopyInto(dst *MonthActive) {
	for d := range 31 {
		for h := range 24 {
			m[d][h].deepCopyInto(&dst[d][h])
		}
	}
}

// getMetric does not verify input.
// Check before calling.
func (m *MonthActive) getMetric(day, hour int8) *MetricActive {
	return &m[day][hour]
}

func (m *MonthActive) WriteTo(b *strings.Builder, withLabel string) {
	for ixDay := range m {
		day := &m[ixDay] // pointer, no copy

		for ixHour := range day {
			m := &day[ixHour] // pointer, no copy

			noRecords := m.RecordsPerPeriod.Load()
			if noRecords == 0 {
				continue
			}

			fmt.Fprintf(
				b,
				"  %-10s day: %02d hour: %02d  records: %d\n",

				withLabel,
				ixDay,
				ixHour,
				noRecords,
			)
		}
	}
}
