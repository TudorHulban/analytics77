package analytics

// Advance is not safe for concurrent callers.
func (r *Registry) Advance() {
	r.muAdvance.Lock()
	r.advance()
	r.muAdvance.Unlock()
}

func (r *Registry) advance() {
	next := (r.CurrentSlot.Load() + 1) % int32(len(r.Slots))

	// new current month slot: just zero the existing MonthActive
	r.zeroSlot(next)

	// move current pointer
	r.CurrentSlot.Store(next)

	nextMonth := r.CalendarMonthCurrentNumber.Load() + 1
	if nextMonth > 12 {
		nextMonth = 1
	}

	r.CalendarMonthCurrentNumber.Store(nextMonth)
}
