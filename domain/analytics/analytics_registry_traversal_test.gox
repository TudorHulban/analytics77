package analytics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPreviousMonthForEach(t *testing.T) {
	r := NewRegistry()

	// Populate some data
	r.GetPreviousMonth()[2][5].RecordsPerPeriod.Store(10)
	r.GetPreviousMonth()[2][7].RecordsPerPeriod.Store(3)
	r.GetPreviousMonth()[10][0].RecordsPerPeriod.Store(1)

	var calls []struct {
		day  int8
		hour int8
		val  uint32
	}

	r.PreviousMonthForEach(
		func(day int8, hour int8, m *MetricActive) {
			calls = append(
				calls,
				struct {
					day  int8
					hour int8
					val  uint32
				}{
					day:  day,
					hour: hour,
					val:  m.RecordsPerPeriod.Load(),
				},
			)
		},
	)

	require.Equal(t,
		3,
		len(calls),
		"expected 3 calls, got %d",
		len(calls),
	)

	require.Equal(t, int8(2), calls[0].day)
	require.Equal(t, int8(5), calls[0].hour)
	require.Equal(t, uint32(10), calls[0].val)

	require.Equal(t, int8(2), calls[1].day)
	require.Equal(t, int8(7), calls[1].hour)
	require.Equal(t, uint32(3), calls[1].val)

	require.Equal(t, int8(10), calls[2].day)
	require.Equal(t, int8(0), calls[2].hour)
	require.Equal(t, uint32(1), calls[2].val)
}

func TestCurrentMonthForEach(t *testing.T) {
	r := NewRegistry()

	// Populate some data
	r.GetCurrentMonth()[0][3].RecordsPerPeriod.Store(5)
	r.GetCurrentMonth()[4][10].RecordsPerPeriod.Store(2)
	r.GetCurrentMonth()[30][23].RecordsPerPeriod.Store(9)

	var calls []struct {
		day  int8
		hour int8
		val  uint32
	}

	r.CurrentMonthForEach(
		func(day, hour int8, m *MetricActive) {
			calls = append(
				calls,
				struct {
					day  int8
					hour int8
					val  uint32
				}{
					day:  day,
					hour: hour,
					val:  m.RecordsPerPeriod.Load(),
				},
			)
		},
	)

	require.Equal(t,
		3,
		len(calls),

		"expected 3 calls, got %d",
		len(calls),
	)

	require.Equal(t, int8(0), calls[0].day)
	require.Equal(t, int8(3), calls[0].hour)
	require.Equal(t, uint32(5), calls[0].val)

	require.Equal(t, int8(4), calls[1].day)
	require.Equal(t, int8(10), calls[1].hour)
	require.Equal(t, uint32(2), calls[1].val)

	require.Equal(t, int8(30), calls[2].day)
	require.Equal(t, int8(23), calls[2].hour)
	require.Equal(t, uint32(9), calls[2].val)
}

func TestHistoryForEach(t *testing.T) {
	r := NewRegistry()

	// Populate some archived data
	r.History[0][1][5].RecordsPerPeriod = 11
	r.History[3][10][0].RecordsPerPeriod = 4
	r.History[6][30][23].RecordsPerPeriod = 99

	type record struct {
		month int8
		day   int8
		hour  int8

		value uint32
	}

	var records []record

	r.HistoryForEach(
		func(month, day, hour int8, m *MetricArchived) {
			records = append(
				records,
				record{
					month: month,
					day:   day,
					hour:  hour,
					value: m.RecordsPerPeriod,
				},
			)
		},
	)

	require.Equal(t,
		3,
		len(records),

		"expected 3 calls, got %d",
		len(records),
	)

	require.Equal(t, int8(0), records[0].month)
	require.Equal(t, int8(1), records[0].day)
	require.Equal(t, int8(5), records[0].hour)
	require.Equal(t, uint32(11), records[0].value)

	require.Equal(t, int8(3), records[1].month)
	require.Equal(t, int8(10), records[1].day)
	require.Equal(t, int8(0), records[1].hour)
	require.Equal(t, uint32(4), records[1].value)

	require.Equal(t, int8(6), records[2].month)
	require.Equal(t, int8(30), records[2].day)
	require.Equal(t, int8(23), records[2].hour)
	require.Equal(t, uint32(99), records[2].value)
}
