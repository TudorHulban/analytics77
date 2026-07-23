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

func (*Registry) MonthForEach(slot *MonthActive, action func(day, hour int8, m *MetricActive)) {
	if action == nil || slot == nil {
		return
	}

	for day := range int8(31) {
		for hour := range int8(24) {
			m := slot.GetMetric(day, hour)

			if m.GetRecordsPerPeriod() != 0 {
				action(day, hour, m)
			}
		}
	}
}

func (r *Registry) HistoryForEach(action func(monthsBack, day, hour int8, m *MetricActive)) {
	if action == nil {
		return
	}

	for months := range uint8(len(r.Slots)) {
		slot, err := r.GetHistorySlot(MonthsBack(months))
		if err != nil {
			continue
		}

		for day := range int8(31) {
			for hour := range int8(24) {
				m := (*slot).GetMetric(day, hour)

				if m.GetRecordsPerPeriod() != 0 {
					action(int8(months), day, hour, m)
				}
			}
		}
	}
}
