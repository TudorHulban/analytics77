package analytics

import "testing"

func BenchmarkPreviousMonthAggregateTopN(b *testing.B) {
	// 1. Prepare registry with synthetic data
	var r Registry

	for day := range int8(31) {
		for hour := range int8(24) {
			m := &r.GetPreviousMonth()[day][hour]
			m.RecordsPerPeriod.Store(1)

			m.TopIPs.Increment("1.1.1.1", uint32(day+1))
			m.TopIPs.Increment("8.8.8.8", uint32(hour+1))

			m.TopCountries.Increment("RO", uint32(day+hour+1))
			m.TopCountries.Increment("US", uint32((day*2)+1))

			m.TopASN.Increment("AS1234", uint32(hour+3))
		}
	}

	// 2. Benchmark
	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		_ = r.PreviousMonthAggregateTopN()
	}
}

func BenchmarkCurrentMonthAggregateTopN(b *testing.B) {
	// 1. Prepare registry with synthetic data
	var r Registry

	for day := range int8(31) {
		for hour := range int8(24) {
			m := &r.GetCurrentMonth()[day][hour]
			m.RecordsPerPeriod.Store(1)

			m.TopIPs.Increment("1.1.1.1", uint32(day+1))
			m.TopIPs.Increment("8.8.8.8", uint32(hour+1))

			m.TopCountries.Increment("RO", uint32(day+hour+1))
			m.TopCountries.Increment("US", uint32((day*2)+1))

			m.TopASN.Increment("AS1234", uint32(hour+3))
		}
	}

	// 2. Benchmark
	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		_ = r.CurrentMonthAggregateTopN()
	}
}

func BenchmarkHistoryAggregateTopN(b *testing.B) {
	// 1. Prepare registry with synthetic data
	var r Registry

	for week := range 7 {
		for day := range int8(31) {
			for hour := range int8(24) {
				m := &r.History[week][day][hour]
				m.RecordsPerPeriod = 1

				m.TopIPs.Names[0] = "1.1.1.1"
				m.TopIPs.Values[0] = uint32(day + hour + 1)

				m.TopCountries.Names[0] = "RO"
				m.TopCountries.Values[0] = uint32((week+1)*3 + int(day))

				m.TopASN.Names[0] = "AS1234"
				m.TopASN.Values[0] = uint32(hour + 5)
			}
		}
	}

	// 2. Benchmark
	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		_ = r.HistoryAggregateTopN()
	}
}
