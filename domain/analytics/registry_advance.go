package analytics

func (r *Registry) Advance() {
	next := (r.CurrentSlot.Load() + 1) % int32(len(r.Slots))

	// new current month slot: just zero the existing MonthActive
	r.zeroSlot(next)

	// move current pointer
	r.CurrentSlot.Store(next)

	r.CalendarMonthCurrentNumber.Add(1)
	r.CalendarMonthCurrentNumber.CompareAndSwap(13, 1)
}
