package analytics

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestZeroSlot(t *testing.T) {
	r := NewRegistry(January)

	ip1 := "1.1.1.1"
	number := 10

	activeSlot := r.GetActiveSlot()

	metric := activeSlot.getMetric(0, 0)

	require.True(t,
		metric.TopIPs.IsZero(),
	)

	valueBefore, errBefore := metric.TopIPs.GetValue("")
	require.ErrorIs(t, errBefore, ErrKeyNotFound)
	require.Zero(t, valueBefore)

	metric.RecordsPerPeriod.Store(1)
	metric.TopIPs.Increment(ip1, uint32(number))

	require.False(t,
		metric.TopIPs.IsZero(),
	)
	require.EqualValues(t,
		number,
		metric.TopIPs.Count(ip1),
	)

	fmt.Println(
		metric.TopIPs.String(),
	)

	r.zeroSlot(r.CurrentSlot.Load())

	require.True(t,
		metric.TopIPs.IsZero(),
	)
	require.Zero(t, metric.TopIPs.Count(ip1))

	require.Zero(t, metric.TopIPs.Count(""))

	valueAfter, errAfter := metric.TopIPs.GetValue("")
	require.ErrorIs(t, errAfter, ErrKeyNotFound)
	require.Zero(t, valueAfter)
}
