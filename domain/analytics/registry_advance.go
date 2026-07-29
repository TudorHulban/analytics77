package analytics

// Advance is not safe for concurrent callers.
func (r *Registry) Advance() {
	r.muAdvance.Lock()
	r.advance()
	r.muAdvance.Unlock()
}

// advance rotates the ring buffer to the next slot.
//
// # ROLLOVER SEMANTICS — TORN READS IN THE OLDEST HISTORY SLOT
//
// zeroSlot runs without holding any reader lock. A concurrent query for the
// oldest history slot (monthsBack == 5) may therefore observe a partially
// zeroed hour: e.g. RecordsPerPeriod already reset to 0 while TopIPs still
// holds stale data, or vice versa.
//
// This is accepted because:
//  1. The window is microseconds-wide — it lasts only for the duration of
//     zeroing a single [31][24]MetricActive array, which happens once per month.
//  2. The affected slot is the *oldest* history slot (6 months back), queried
//     least frequently in practice.
//  3. All individual field accesses are atomic, so this is a torn read, not a
//     Go data race.
//  4. Preventing it would require either a reader-side version check (adds
//     overhead to every query) or a pointer swap (breaks the zero-allocation,
//     fixed-address invariants of this package).
//
// Callers that need strictly consistent snapshots for the oldest history slot
// must serialize with Advance() externally.
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
