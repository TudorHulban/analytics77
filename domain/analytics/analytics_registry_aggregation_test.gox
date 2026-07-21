package analytics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPreviousMonthTotalRecords(t *testing.T) {
	r := NewRegistry()

	// Populate some data
	r.GetPreviousMonth()[0][0].RecordsPerPeriod.Store(5)
	r.GetPreviousMonth()[0][1].RecordsPerPeriod.Store(3)
	r.GetPreviousMonth()[10][5].RecordsPerPeriod.Store(12)
	r.GetPreviousMonth()[30][23].RecordsPerPeriod.Store(20)

	got := r.PreviousMonthTotalRecords()

	want := int32(5 + 3 + 12 + 20)

	require.EqualValues(t,
		want,
		got,

		"total mismatch: got %d want %d",
		got,
		want,
	)
}

func TestPreviousMonthTotalRecordsForDay(t *testing.T) {
	r := NewRegistry()

	// day 12: hours 0, 3, 7 have data
	r.GetPreviousMonth()[12][0].RecordsPerPeriod.Store(5)
	r.GetPreviousMonth()[12][3].RecordsPerPeriod.Store(8)
	r.GetPreviousMonth()[12][7].RecordsPerPeriod.Store(20)

	got := r.PreviousMonthTotalRecordsForDay(12)

	want := int32(5 + 8 + 20)

	require.EqualValues(t,
		want,
		got,

		"total mismatch: got %d want %d",
		got,
		want,
	)
}

func TestCurrentMonthTotalRecords(t *testing.T) {
	r := NewRegistry()

	// Populate some data
	r.GetCurrentMonth()[0][0].RecordsPerPeriod.Store(7)
	r.GetCurrentMonth()[5][12].RecordsPerPeriod.Store(13)
	r.GetCurrentMonth()[30][23].RecordsPerPeriod.Store(2)

	got := r.CurrentMonthTotalRecords()

	want := int32(7 + 13 + 2)

	require.EqualValues(t,
		want,
		got,

		"total mismatch: got %d want %d",
		got,
		want,
	)
}

func TestCurrentMonthTotalRecordsForDay(t *testing.T) {
	r := NewRegistry()

	// day 8: hours 2, 11, 20 have data
	r.GetCurrentMonth()[8][2].RecordsPerPeriod.Store(4)
	r.GetCurrentMonth()[8][11].RecordsPerPeriod.Store(9)
	r.GetCurrentMonth()[8][20].RecordsPerPeriod.Store(16)

	got := r.CurrentMonthTotalRecordsForDay(8)

	want := int32(4 + 9 + 16)

	require.EqualValues(t,
		want,
		got,

		"total mismatch: got %d want %d",
		got,
		want,
	)
}

func TestPreviousMonthAggregateTopN(t *testing.T) {
	r := NewRegistry()

	// Bucket 1
	r.GetPreviousMonth()[0][0].RecordsPerPeriod.Store(1)
	r.GetPreviousMonth()[0][0].TopIPs.Increment("1.1.1.1", 3)
	r.GetPreviousMonth()[0][0].TopCountries.Increment("RO", 2)

	// Bucket 2
	r.GetPreviousMonth()[5][12].RecordsPerPeriod.Store(1)
	r.GetPreviousMonth()[5][12].TopIPs.Increment("1.1.1.1", 7)
	r.GetPreviousMonth()[5][12].TopIPs.Increment("8.8.8.8", 4)
	r.GetPreviousMonth()[5][12].TopCountries.Increment("US", 5)

	// Bucket 3
	r.GetPreviousMonth()[30][23].RecordsPerPeriod.Store(1)
	r.GetPreviousMonth()[30][23].TopCountries.Increment("RO", 1)
	r.GetPreviousMonth()[30][23].TopASN.Increment("AS1234", 9)

	agg := r.PreviousMonthAggregateTopN()

	// IPs
	require.Equal(t, uint32(10), agg.IPs.Count("1.1.1.1"))
	require.Equal(t, uint32(4), agg.IPs.Count("8.8.8.8"))

	// Countries
	require.Equal(t, uint32(3), agg.Countries.Count("RO"))
	require.Equal(t, uint32(5), agg.Countries.Count("US"))

	// ASN
	require.Equal(t, uint32(9), agg.ASN.Count("AS1234"))
}

func TestCurrentMonthAggregateTopN(t *testing.T) {
	r := NewRegistry()

	// Bucket 1
	r.GetCurrentMonth()[0][0].RecordsPerPeriod.Store(1)
	r.GetCurrentMonth()[0][0].TopIPs.Increment("1.1.1.1", 3)
	r.GetCurrentMonth()[0][0].TopCountries.Increment("RO", 2)

	// Bucket 2
	r.GetCurrentMonth()[5][12].RecordsPerPeriod.Store(1)
	r.GetCurrentMonth()[5][12].TopIPs.Increment("1.1.1.1", 7)
	r.GetCurrentMonth()[5][12].TopIPs.Increment("8.8.8.8", 4)
	r.GetCurrentMonth()[5][12].TopCountries.Increment("US", 5)

	// Bucket 3
	r.GetCurrentMonth()[30][23].RecordsPerPeriod.Store(1)
	r.GetCurrentMonth()[30][23].TopCountries.Increment("RO", 1)
	r.GetCurrentMonth()[30][23].TopASN.Increment("AS1234", 9)

	agg := r.CurrentMonthAggregateTopN()

	// IPs
	require.Equal(t, uint32(10), agg.IPs.Count("1.1.1.1"))
	require.Equal(t, uint32(4), agg.IPs.Count("8.8.8.8"))

	// Countries
	require.Equal(t, uint32(3), agg.Countries.Count("RO"))
	require.Equal(t, uint32(5), agg.Countries.Count("US"))

	// ASN
	require.Equal(t, uint32(9), agg.ASN.Count("AS1234"))
}

func TestHistoryAggregateTopN(t *testing.T) {
	r := NewRegistry()

	// Month 0
	r.History[0][1][5].RecordsPerPeriod = 1
	r.History[0][1][5].TopIPs.Names[0] = "1.1.1.1"
	r.History[0][1][5].TopIPs.Values[0] = 3

	// Month 3
	r.History[3][10][0].RecordsPerPeriod = 1
	r.History[3][10][0].TopIPs.Names[0] = "1.1.1.1"
	r.History[3][10][0].TopIPs.Values[0] = 4
	r.History[3][10][0].TopCountries.Names[0] = "US"
	r.History[3][10][0].TopCountries.Values[0] = 5

	// Month 6
	r.History[6][30][23].RecordsPerPeriod = 1
	r.History[6][30][23].TopCountries.Names[0] = "RO"
	r.History[6][30][23].TopCountries.Values[0] = 1
	r.History[6][30][23].TopASN.Names[0] = "AS1234"
	r.History[6][30][23].TopASN.Values[0] = 9

	agg := r.HistoryAggregateTopN()

	// IPs
	require.Equal(t, uint32(7), agg.IPs.Count("1.1.1.1"))

	// Countries
	require.Equal(t, uint32(1), agg.Countries.Count("RO"))
	require.Equal(t, uint32(5), agg.Countries.Count("US"))

	// ASN
	require.Equal(t, uint32(9), agg.ASN.Count("AS1234"))
}
