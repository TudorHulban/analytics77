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

// HistoryForEach iterates over all *physical* history slots except the current slot.
// For each non‑zero metric, it invokes `action` with:
//
//	monthsBack = logical distance from the current slot (0..6)
//	day        = day index (0..30)
//	hour       = hour index (0..23)
//	m          = pointer to the actual MetricActive stored in the physical slot
//
// Note:
//   - `monthsBack` is computed relative to the current slot via modulo arithmetic.
//   - The physical slot index and the logical monthsBack value are not the same.
//   - The callback receives the real metric from the underlying slot, not a remapped copy.
func (r *Registry) HistoryForEach(action func(monthsBack, day, hour int8, m *MetricActive)) {
	if action == nil {
		return
	}

	currentSlot := uint8(r.CurrentSlot.Load()) //nolint:gosec

	for monthIx := uint8(0); monthIx < uint8(len(r.Slots)); monthIx++ {
		if monthIx == currentSlot {
			continue
		}

		// Compute logical monthsBack relative to current slot
		monthsBack := int8((currentSlot - monthIx + 7) % 7)
		slot := r.Slots[monthIx]

		for day := int8(0); day < 31; day++ {
			for hour := int8(0); hour < 24; hour++ {
				m := &(*slot)[day][hour]

				if m.GetRecordsPerPeriod() != 0 {
					action(monthsBack, day, hour, m)
				}
			}
		}
	}
}
