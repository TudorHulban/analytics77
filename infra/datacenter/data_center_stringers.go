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

func registryString(r *analytics.Registry, b *strings.Builder) {
	monthActiveString("current", r.GetActiveSlot()[:], b)
	monthActiveString("previous", r.GetPreviousSlot()[:], b)

	// archived months (the remaining 5 slots)
	for ix := range 7 {
		month := &r.Slots[ix]

		// skip current and previous
		if ix == int(r.CurrentSlot.Load()) {
			continue
		}

		if ix == int((r.CurrentSlot.Load()+7-1)%7) {
			continue
		}

		monthActiveString(
			fmt.Sprintf("slot[%d]", ix),
			month[:],
			b,
		)
	}
}
