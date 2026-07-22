package analytics

// Snapshot writes a textual snapshot of the registry into w.
// func (r *Registry) Snapshot(w io.Writer) error {
// 	currentMonth := r.GetCurrentMonth()

// 	for day := range 31 {
// 		for hour := range 24 {
// 			m := &currentMonth[day][hour]

// 			if m.RecordsPerPeriod.Load() == 0 {
// 				continue
// 			}

// 			if _, errPrintCurrent := fmt.Fprintf(
// 				w,

// 				"current day:'%02d' hour:'%02d' records:'%d'\n",
// 				day,
// 				hour,
// 				m.RecordsPerPeriod.Load(),
// 			); errPrintCurrent != nil {
// 				return errPrintCurrent
// 			}
// 		}
// 	}

// 	previousMonth := r.GetPreviousMonth()

// 	for day := range 31 {
// 		for hour := range 24 {
// 			m := &previousMonth[day][hour]
// 			if m.RecordsPerPeriod.Load() == 0 {
// 				continue
// 			}

// 			if _, errPrintPrevious := fmt.Fprintf(
// 				w,

// 				"previous day: %02d hour: %02d records:%d\n",
// 				day,
// 				hour,
// 				m.RecordsPerPeriod.Load(),
// 			); errPrintPrevious != nil {
// 				return errPrintPrevious
// 			}
// 		}
// 	}

// 	for h := range 7 {
// 		for day := range 31 {
// 			for hour := range 24 {
// 				m := &r.History[h][day][hour]
// 				if m.RecordsPerPeriod == 0 {
// 					continue
// 				}

// 				if _, errPrintHistory := fmt.Fprintf(
// 					w,

// 					"history[%d] day: %02d hour: %02d records:%d\n",
// 					h,
// 					day,
// 					hour,
// 					m.RecordsPerPeriod,
// 				); errPrintHistory != nil {
// 					return errPrintHistory
// 				}
// 			}
// 		}
// 	}

// 	return nil
// }
