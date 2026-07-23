package analytics

import (
	"fmt"
	"io"
)

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

			if _, err := fmt.Fprintf(
				w,
				"current day:%02d hour:%02d records:%d\n",
				day,
				hour,
				rec,
			); err != nil {
				return err
			}
		}
	}

	// --- history months (including previous) ---
	for months := range uint8(len(r.Slots)) {
		slot, err := r.GetHistorySlot(MonthsBack(months))
		if err != nil {
			continue
		}

		label := fmt.Sprintf("history[%d]", months)

		for day := range 31 {
			for hour := range 24 {
				m := &slot[day][hour]

				rec := m.RecordsPerPeriod.Load()
				if rec == 0 {
					continue
				}

				if _, err := fmt.Fprintf(
					w,
					"%s day:%02d hour:%02d records:%d\n",
					label,
					day,
					hour,
					rec,
				); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
