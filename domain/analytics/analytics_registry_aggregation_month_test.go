package analytics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHistoryAggregateTopNForMonth(t *testing.T) {
	r := NewRegistry(0)

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

	slot := r.Slots[0]

	// Day 0, Hour 0
	m00 := slot.GetMetric(0, 0)
	m00.RecordsPerPeriod.Store(1)
	m00.TopIPs.Increment(ip1, 10)
	m00.TopIPs.Increment(ip2, 5)

	// Day 0, Hour 1
	m01 := slot.GetMetric(0, 1)
	m01.RecordsPerPeriod.Store(1)
	m01.TopIPs.Increment(ip1, 3)

	// Day 1, Hour 0
	m10 := slot.GetMetric(1, 0)
	m10.RecordsPerPeriod.Store(1)
	m10.TopIPs.Increment(ip2, 7)

	// Day 2, Hour 5
	m25 := slot.GetMetric(2, 5)
	m25.RecordsPerPeriod.Store(1)
	m25.TopIPs.Increment(ip3, 11)

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

	require.EqualValues(t, 13, got1, "aggregated value for IP 1.1.1.1 incorrect")
	require.EqualValues(t, 12, got2, "aggregated value for IP 2.2.2.2 incorrect")
	require.EqualValues(t, 11, got3, "aggregated value for IP 3.3.3.3 incorrect")
}
