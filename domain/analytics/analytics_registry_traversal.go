package analytics

func (r *Registry) PreviousMonthForEach(action func(day, hour int8, m *MetricActive)) {
	if action == nil {
		return
	}

	previousSlot := r.GetPreviousSlot()

	for day := range int8(31) {
		for hour := range int8(24) {
			m := (*previousSlot).GetMetric(day, hour)

			if m.GetRecordsPerPeriod() != 0 {
				action(day, hour, m)
			}
		}
	}
}

func (r *Registry) CurrentMonthForEach(action func(day, hour int8, m *MetricActive)) {
	if action == nil {
		return
	}

	activeSlot := r.GetActiveSlot()

	for day := range int8(31) {
		for hour := range int8(24) {
			m := (*activeSlot).GetMetric(day, hour)

			if m.GetRecordsPerPeriod() != 0 {
				action(day, hour, m)
			}
		}
	}
}

func (r *Registry) HistoryForEach(action func(monthIndex, day, hour int8, m *MetricActive)) {
	if action == nil {
		return
	}

	for month := range int8(7) {
		slot := &r.Slots[month]

		for day := range int8(31) {
			for hour := range int8(24) {
				m := (*slot).GetMetric(day, hour)

				if m.GetRecordsPerPeriod() != 0 {
					action(month, day, hour, m)
				}
			}
		}
	}
}
