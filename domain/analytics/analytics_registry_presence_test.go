package analytics

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/analytics77/domain/dhelpers"
	"github.com/tudorhulban/hxhelpers"
)

func TestPreviousMonthDaysWithData(t *testing.T) {
	r := NewRegistry()

	// day 0: empty
	// day 1: hour 5 has data
	r.GetPreviousSlot()[1][5].RecordsPerPeriod.Store(10)

	// day 7: hour 0 has data
	r.GetPreviousSlot()[7][0].RecordsPerPeriod.Store(1)

	daysWithData, howMany := r.PreviousMonthDaysWithData()

	want := []int8{1, 7}

	require.EqualValues(t,
		howMany,
		len(want),

		"len mismatch: got %d want %d",
		howMany,
		len(want),
	)

	require.Empty(t,
		hxhelpers.NotInSliceSource(daysWithData[:howMany], want...),
	)
}

func TestPreviousMonthHoursWithData(t *testing.T) {
	r := NewRegistry()

	const dayStorage int8 = 3

	// day 3: hours 0, 5, 23 have data
	r.GetPreviousSlot()[dayStorage][0].RecordsPerPeriod.Store(1)
	r.GetPreviousSlot()[dayStorage][5].RecordsPerPeriod.Store(10)
	r.GetPreviousSlot()[dayStorage][23].RecordsPerPeriod.Store(7)

	daysWithData, howMany := r.PreviousMonthHoursWithData(dhelpers.IndexToCalendarDay(dayStorage))

	want := []int8{0, 5, 23}

	require.EqualValues(t,
		howMany,
		len(want),

		"len mismatch: got %d want %d",
		howMany,
		len(want),
	)

	require.Empty(t,
		hxhelpers.NotInSliceSource(daysWithData[:howMany], want...),
	)
}

func TestCurrentMonthDaysWithData(t *testing.T) {
	r := NewRegistry()

	// day 0: empty
	// day 2: hour 10 has data
	r.GetActiveSlot()[2][10].RecordsPerPeriod.Store(5)

	// day 15: hour 0 has data
	r.GetActiveSlot()[15][0].RecordsPerPeriod.Store(1)

	daysWithData, howMany := r.CurrentMonthDaysWithData()

	want := []int8{2, 15}

	require.EqualValues(t,
		howMany,
		len(want),

		"len mismatch: got %d want %d",
		howMany,
		len(want),
	)

	require.Empty(t,
		hxhelpers.NotInSliceSource(daysWithData[:howMany], want...),
	)
}

func TestCurrentMonthHoursWithData(t *testing.T) {
	r := NewRegistry()

	storageDay := 4

	// day 4: hours 1, 7, 22 have data
	r.GetActiveSlot()[storageDay][1].RecordsPerPeriod.Store(3)
	r.GetActiveSlot()[storageDay][7].RecordsPerPeriod.Store(9)
	r.GetActiveSlot()[storageDay][22].RecordsPerPeriod.Store(1)

	daysWithData, howMany := r.
		CurrentMonthHoursWithData(dhelpers.IndexToCalendarDay(int8(storageDay)))

	want := []int8{1, 7, 22}

	require.EqualValues(t,
		howMany,
		len(want),

		"len mismatch: got %d want %d",
		howMany,
		len(want),
	)

	require.Empty(t,
		hxhelpers.NotInSliceSource(daysWithData[:howMany], want...),
	)
}
