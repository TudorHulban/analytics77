package analytics

import "fmt"

type MonthsBack uint8

var (
	FromPreviousMonth    MonthsBack = 0 //nolint:revive
	FromTwoMonthsAgo     MonthsBack = 1
	FromThreeMonthsAgo   MonthsBack = 2
	FromFourMonthsAgo    MonthsBack = 3
	FromFiveMonthsAgo    MonthsBack = 4
	FromLastHistoryMonth MonthsBack = 5
)

func (m MonthsBack) String() string {
	switch m {
	case FromPreviousMonth:
		return "previous month"
	case FromTwoMonthsAgo:
		return "two months ago"
	case FromThreeMonthsAgo:
		return "three months ago"
	case FromFourMonthsAgo:
		return "four months ago"
	case FromFiveMonthsAgo:
		return "five months ago"
	case FromLastHistoryMonth:
		return "six months ago"

	default:
		return fmt.Sprintf("monthsBack(%d)", uint8(m))
	}
}

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
	if monthsBack >= MonthsBack(len(r.Slots)-1) {
		return nil,
			ErrInvalidInput
	}

	prev := (r.CurrentSlot.Load() + int32(len(r.Slots)) - int32(monthsBack+1)) % int32(len(r.Slots))

	return r.Slots[prev], nil
}
