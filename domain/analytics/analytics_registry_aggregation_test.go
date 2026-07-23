package analytics

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPreviousMonthTotalRecords(t *testing.T) {
	r := NewRegistry()

	previousSlot := r.GetPreviousSlot()

	// Populate some data
	previousSlot[0][0].RecordsPerPeriod.Store(5)
	previousSlot[0][1].RecordsPerPeriod.Store(3)
	previousSlot[10][5].RecordsPerPeriod.Store(12)
	previousSlot[30][23].RecordsPerPeriod.Store(20)

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

	previousSlot := r.GetPreviousSlot()

	// Populate some data
	// day 12: hours 0, 3, 7 have data
	previousSlot[12][0].RecordsPerPeriod.Store(5)
	previousSlot[12][3].RecordsPerPeriod.Store(8)
	previousSlot[12][7].RecordsPerPeriod.Store(20)

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

	activeSlot := r.GetActiveSlot()

	// Populate some data
	activeSlot[0][0].RecordsPerPeriod.Store(7)
	activeSlot[5][12].RecordsPerPeriod.Store(13)
	activeSlot[30][23].RecordsPerPeriod.Store(2)

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

	activeSlot := r.GetActiveSlot()

	// Populate some data for
	// day 8: hours 2, 11, 20 have data
	activeSlot[8][2].RecordsPerPeriod.Store(4)
	activeSlot[8][11].RecordsPerPeriod.Store(9)
	activeSlot[8][20].RecordsPerPeriod.Store(16)

	got := r.CurrentMonthTotalRecordsForDay(8)

	require.Zero(t, r.CurrentMonthTotalRecordsForDay(7))

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

	previousSlot := r.GetPreviousSlot()

	// Bucket 1
	previousSlot[0][0].RecordsPerPeriod.Store(1)
	previousSlot[0][0].TopIPs.Increment("1.1.1.1", 3)
	previousSlot[0][0].TopCountries.Increment("RO", 2)

	// Bucket 2
	previousSlot[5][12].RecordsPerPeriod.Store(1)
	previousSlot[5][12].TopIPs.Increment("1.1.1.1", 7)
	previousSlot[5][12].TopIPs.Increment("8.8.8.8", 4)
	previousSlot[5][12].TopCountries.Increment("US", 5)

	// Bucket 3
	previousSlot[30][23].RecordsPerPeriod.Store(1)
	previousSlot[30][23].TopCountries.Increment("RO", 1)
	previousSlot[30][23].TopASN.Increment("AS1234", 9)

	agg := r.PreviousMonthAggregateTopN()

	fmt.Println(agg.String())

	// IPs
	require.Equal(t,
		uint32(10),
		agg.IPs.Count("1.1.1.1"),
	)
	require.Equal(t,
		uint32(4),
		agg.IPs.Count("8.8.8.8"),
	)

	// Countries
	require.Equal(t, uint32(3), agg.Countries.Count("RO"))
	require.Equal(t, uint32(5), agg.Countries.Count("US"))

	// ASN
	require.Equal(t, uint32(9), agg.ASN.Count("AS1234"))
}

func TestCurrentMonthAggregateTopN(t *testing.T) {
	r := NewRegistry()

	// Bucket 1
	r.GetActiveSlot()[0][0].RecordsPerPeriod.Store(1)
	r.GetActiveSlot()[0][0].TopIPs.Increment("1.1.1.1", 3)
	r.GetActiveSlot()[0][0].TopCountries.Increment("RO", 2)

	// Bucket 2
	r.GetActiveSlot()[5][12].RecordsPerPeriod.Store(1)
	r.GetActiveSlot()[5][12].TopIPs.Increment("1.1.1.1", 7)
	r.GetActiveSlot()[5][12].TopIPs.Increment("8.8.8.8", 4)
	r.GetActiveSlot()[5][12].TopCountries.Increment("US", 5)

	// Bucket 3
	r.GetActiveSlot()[30][23].RecordsPerPeriod.Store(1)
	r.GetActiveSlot()[30][23].TopCountries.Increment("RO", 1)
	r.GetActiveSlot()[30][23].TopASN.Increment("AS1234", 9)

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

	_, errHistory7 := r.GetHistorySlot(7)
	require.Error(t, errHistory7)

	// Month 0
	slot0, errHistory0 := r.GetHistorySlot(0)
	require.NoError(t, errHistory0)

	ip := "1.1.1.1"
	countryRO := "RO"
	countryUS := "US"

	slot0[1][5].RecordsPerPeriod.Store(1)
	slot0[1][5].TopIPs.Names[0].Store(&ip)
	slot0[1][5].TopIPs.Values[0].Store(3)

	// Month 3
	slot3, errHistory3 := r.GetHistorySlot(3)
	require.NoError(t, errHistory3)

	slot3[10][0].RecordsPerPeriod.Store(1)
	slot3[10][0].TopIPs.Names[0].Store(&ip)
	slot3[10][0].TopIPs.Values[0].Store(4)
	slot3[10][0].TopCountries.Names[0].Store(&countryUS)
	slot3[10][0].TopCountries.Values[0].Store(5)

	// Month 6
	slot6, errHistory6 := r.GetHistorySlot(6)
	require.NoError(t, errHistory6)

	asn := "AS1234"

	slot6[30][23].RecordsPerPeriod.Store(1)
	slot6[30][23].TopCountries.Names[0].Store(&countryRO)
	slot6[30][23].TopCountries.Values[0].Store(1)
	slot6[30][23].TopASN.Names[0].Store(&asn)
	slot6[30][23].TopASN.Values[0].Store(9)

	agg := r.HistoryAggregateTopN()

	// IPs
	require.Equal(t,
		uint32(7),
		agg.IPs.Count(ip),
	)

	// Countries
	require.Equal(t,
		uint32(1),
		agg.Countries.Count(countryRO),
	)
	require.Equal(t,
		uint32(5),
		agg.Countries.Count(countryUS),
	)

	// ASN
	require.Equal(t,
		uint32(9),
		agg.ASN.Count(asn),
	)
}
