package analytics

type MonthsBack uint8

var (
	FromPreviousMonth    MonthsBack = 0 //nolint:revive
	FromTwoMonthsAgo     MonthsBack = 1
	FromThreeMonthsAgo   MonthsBack = 2
	FromFourMonthsAgo    MonthsBack = 3
	FromFiveMonthsAgo    MonthsBack = 4
	FromLastHistoryMonth MonthsBack = 5
)

// monthsBack = 0 → previous month
//
// monthsBack = 1 → two months ago
//
// monthsBack = 2 → three months ago
//
// monthsBack = 3 → four months ago
//
// monthsBack = 4 → five months ago
//
// monthsBack = 5 → six months ago
func (r *Registry) GetHistorySlot(monthsBack MonthsBack) (*MonthActive, error) {
	if monthsBack >= MonthsBack(len(r.Slots)) {
		return nil,
			ErrInvalidInput
	}

	prev := (r.CurrentSlot.Load() + int32(len(r.Slots)) - int32(monthsBack+1)) % int32(len(r.Slots))

	return &r.Slots[prev], nil
}
