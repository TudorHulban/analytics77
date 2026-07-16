package analytics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRollover_HistoryShift(t *testing.T) {
	r := NewRegistry()

	// Seed each history slot with a unique record count so we can track them.
	// History[0]=100, History[1]=200, ..., History[6]=700
	for ix := range 7 {
		r.History[ix][0][0].
			RecordsPerPeriod = uint32((ix + 1) * 100)
	}

	// Seed MonthPrevious so we know what should land in History[0]
	r.GetPreviousMonth()[0][0].RecordsPerPeriod.Store(999)

	// ── Act ──
	r.Rollover()

	// ── Assert ──
	want := [7]uint32{
		999, // History[0] = came from MonthPrevious
		100, // History[1] = old History[0]
		200, // History[2] = old History[1]
		300, // History[3] = old History[2]
		400, // History[4] = old History[3]
		500, // History[5] = old History[4]
		600, // History[6] = old History[5] (700 should be dropped)
	}

	for ix, expected := range want {
		got := r.History[ix][0][0].RecordsPerPeriod

		require.Equal(t,
			expected,
			got,

			"History[%d]: got %d, want %d",
			ix,
			got,
			expected,
		)
	}
}
