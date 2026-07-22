package analytics

func (r *Registry) PreviousMonthForEach(action func(day, hour int8, m *MetricActive)) {
	if action == nil {
		return
	}

	previousSlot := r.GetPreviousSlot()

	for day := int8(0); day < 31; day++ {
		for hour := int8(0); hour < 24; hour++ {
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

	for month := int8(0); month < 7; month++ {
		slot := &r.Slots[month]

		for day := int8(0); day < 31; day++ {
			for hour := int8(0); hour < 24; hour++ {
				m := (*slot).GetMetric(day, hour)

				if m.GetRecordsPerPeriod() != 0 {
					action(month, day, hour, m)
				}
			}
		}
	}
}
