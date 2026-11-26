package analytics

import (
	"fmt"
	"io"
	"strings"
)

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

// Snapshot writes a textual snapshot of the registry into w.
func (r *Registry) Snapshot(w io.Writer) error {
	currentMonth := r.GetActiveSlot()

	for day := range 31 {
		for hour := range 24 {
			m := &currentMonth[day][hour]

			rec := m.RecordsPerPeriod.Load()
			if rec == 0 {
				continue
			}

			if _, errPrint := fmt.Fprintf(
				w,

				"current day: '%02d' hour: '%02d' records: '%d'\n",
				day,
				hour,
				rec,
			); errPrint != nil {
				return errPrint
			}
		}
	}

	for back := MonthsBack(0); back <= FromLastHistoryMonth; back++ {
		slot, err := r.GetHistorySlot(back)
		if err != nil {
			continue
		}

		label := fmt.Sprintf("history[%s]", back.String())

		for day := range 31 {
			for hour := range 24 {
				m := &slot[day][hour]

				rec := m.RecordsPerPeriod.Load()
				if rec == 0 {
					continue
				}

				if _, errPrint := fmt.Fprintf(
					w,

					"%s day: '%02d' hour: '%02d' records: '%d'\n",
					label,
					day,
					hour,
					rec,
				); errPrint != nil {
					return errPrint
				}
			}
		}
	}

	return nil
}
