package datacenter

import (
	"fmt"
	"strings"

	"github.com/tudorhulban/analytics77/domain/analytics"
)

func monthActiveString(label string, month []analytics.DayActive, b *strings.Builder) {
	for ixDay := range month {
		day := &month[ixDay] // pointer, no copy

		for ixHour := range day {
			m := &day[ixHour] // pointer, no copy

			noRecords := m.RecordsPerPeriod.Load()
			if noRecords == 0 {
				continue
			}

			fmt.Fprintf(
				b,
				"  %-10s day%02d hour%02d  records:%d\n",

				label,
				ixDay,
				ixHour,
				noRecords,
			)
		}
	}
}

func monthArchivedString(label string, month []analytics.DayArchived, b *strings.Builder) {
	for ixDay := range month {
		day := &month[ixDay] // pointer, no copy

		for ixHour := range day {
			m := &day[ixHour] // pointer, no copy

			noRecords := m.RecordsPerPeriod
			if noRecords == 0 {
				continue
			}

			fmt.Fprintf(
				b,
				"  %-10s day%02d hour%02d  records:%d\n",

				label,
				ixDay,
				ixHour,
				noRecords,
			)
		}
	}
}

func registryString(r *analytics.Registry, b *strings.Builder) {
	monthActiveString("current", r.GetCurrentMonth()[:], b)
	monthActiveString("previous", r.GetPreviousMonth()[:], b)

	for ix, month := range r.History {
		monthArchivedString(
			fmt.Sprintf("history[%d]", ix),
			month[:],
			b,
		)
	}
}
