package analytics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHistoryAggregateTopNForMonth(t *testing.T) {
	r := NewRegistry(January)

	ip1 := "1.1.1.1"
	ip2 := "2.2.2.2"
	ip3 := "3.3.3.3"

	// Pattern:
	//   Day 0, Hour 0 → ip1 = 10, ip2 = 5
	//   Day 0, Hour 1 → ip1 = 3
	//   Day 1, Hour 0 → ip2 = 7
	//   Day 2, Hour 5 → ip3 = 11
	//
	// Expected final aggregation:
	//   ip1 = 13
	//   ip2 = 12
	//   ip3 = 11

	activeSlot := r.GetActiveSlot()

	// Day 0, Hour 0
	metric00 := activeSlot.getMetric(0, 0)
	metric00.RecordsPerPeriod.Store(2)
	metric00.TopIPs.Increment(ip1, 10)
	metric00.TopIPs.Increment(ip2, 5)

	// Day 0, Hour 1
	metric01 := activeSlot.getMetric(0, 1)
	metric01.RecordsPerPeriod.Store(1)
	metric01.TopIPs.Increment(ip1, 3)

	// Day 1, Hour 0
	metric10 := activeSlot.getMetric(1, 0)
	metric10.RecordsPerPeriod.Store(1)
	metric10.TopIPs.Increment(ip2, 7)

	// Day 2, Hour 5
	metric25 := activeSlot.getMetric(2, 5)
	metric25.RecordsPerPeriod.Store(1)
	metric25.TopIPs.Increment(ip3, 11)

	r.Advance() // moves to previous month the slot we used for insert.

	// Aggregate
	agg, err := r.HistoryAggregateTopNForMonth(0)
	require.NoError(t, err)

	var (
		got1 uint32
		got2 uint32
		got3 uint32
	)

	// Extract final aggregated values from TopIPs
	for ix := range 7 {
		namePtr := agg.IPs.Names[ix].Load()
		if namePtr == nil {
			continue
		}

		switch *namePtr {
		case ip1:
			got1 = agg.IPs.Values[ix].Load()
		case ip2:
			got2 = agg.IPs.Values[ix].Load()
		case ip3:
			got3 = agg.IPs.Values[ix].Load()
		}
	}

	require.EqualValues(t,
		10+3,
		got1,

		"aggregated value for IP 1.1.1.1 incorrect, got:%d",
		got1,
	)
	require.EqualValues(t,
		12,
		got2,

		"aggregated value for IP 2.2.2.2 incorrect, got:%d",
		got2,
	)
	require.EqualValues(t,
		11,
		got3,

		"aggregated value for IP 3.3.3.3 incorrect, got:%d",
		got3,
	)
}
