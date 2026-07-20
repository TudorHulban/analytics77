package analytics

type ActionActive func(day int8, hour int8, m *MetricActive)

func (r *Registry) PreviousMonthForEach(action ActionActive) {
	if action == nil {
		return
	}

	previousMonth := r.GetPreviousMonth()

	for day := range int8(31) {
		for hour := range int8(24) {
			m := &previousMonth[day][hour]

			if m.RecordsPerPeriod.Load() != 0 {
				action(day, hour, m)
			}
		}
	}
}

func (r *Registry) CurrentMonthForEach(action ActionActive) {
	if action == nil {
		return
	}

	currentMonth := r.GetCurrentMonth()

	for day := range int8(31) {
		for hour := range int8(24) {
			m := &currentMonth[day][hour]

			if m.RecordsPerPeriod.Load() != 0 {
				action(day, hour, m)
			}
		}
	}
}

type ActionArchive func(month, day, hour int8, m *MetricArchived)

func (r *Registry) HistoryForEach(action ActionArchive) {
	if action == nil {
		return
	}

	for month := range int8(7) {
		for day := range int8(31) {
			for hour := range int8(24) {
				m := &r.History[month][day][hour]

				if m.RecordsPerPeriod != 0 {
					action(month, day, hour, m)
				}
			}
		}
	}
}
