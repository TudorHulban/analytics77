package analytics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdvanceAndHistoryShift(t *testing.T) {
	r := NewRegistry(0)

	// Seed each slot with a unique record count so we can track them.
	for ix := range 7 {
		r.Slots[ix][0][0].RecordsPerPeriod.Store(uint32((ix + 1) * 100))
	}

	oldCurrent := r.CurrentSlot.Load()

	r.Advance()

	newCurrent := r.CurrentSlot.Load()

	// newCurrent must be oldCurrent+1 mod 7
	require.Equal(t,
		(int32(oldCurrent)+1)%7,
		newCurrent,
		"CurrentSlot did not advance correctly",
	)

	// new current slot must be zeroed
	require.Equal(t,
		uint32(0),
		r.Slots[newCurrent][0][0].RecordsPerPeriod.Load(),
		"new current slot must be zeroed",
	)

	// all other slots must remain unchanged
	for ix := range 7 {
		if ix == int(newCurrent) {
			continue
		}

		expected := uint32((ix + 1) * 100)
		got := r.Slots[ix][0][0].RecordsPerPeriod.Load()

		require.Equal(t,
			expected,
			got,

			"slot[%d] changed unexpectedly",
			ix,
		)
	}
}
