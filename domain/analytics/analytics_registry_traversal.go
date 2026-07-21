package analytics

func (r *Registry[T, Data]) PreviousMonthForEach(action func(day, hour int8, m *TData)) {
	if action == nil {
		return
	}

	previousSlot := r.Ring.GetPreviousSlot()

	for day := int8(0); day < 31; day++ {
		for hour := int8(0); hour < 24; hour++ {
			m := (*previousSlot).GetMetric(day, hour)

			if (*m).GetRecordsPerPeriod() != 0 {
				action(day, hour, m)
			}
		}
	}
}

func (r *Registry[T, Data]) CurrentMonthForEach(action func(day, hour int8, m *TData)) {
	if action == nil {
		return
	}

	currentSlot := r.GetCurrentMonth()

	for day := range int8(31) {
		for hour := range int8(24) {
			m := (*currentSlot).GetMetric(day, hour)

			if (*m).GetRecordsPerPeriod() != 0 {
				action(day, hour, m)
			}
		}
	}
}

func (r *Registry[T, Data]) HistoryForEach(action func(monthIndex, day, hour int8, m *TData)) {
	if action == nil {
		return
	}

	for month := int8(0); month < 7; month++ {
		slot := &r.Ring.slots[month]

		for day := int8(0); day < 31; day++ {
			for hour := int8(0); hour < 24; hour++ {
				m := (*slot).GetMetric(day, hour)

				if (*m).GetRecordsPerPeriod() != 0 {
					action(month, day, hour, m)
				}
			}
		}
	}
}
