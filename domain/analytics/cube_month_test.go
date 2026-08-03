package analytics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMonthActiveDeepCopyInto(t *testing.T) {
	var src MonthActive

	// Seed some data across different days/hours
	src[0][0].RecordsPerPeriod.Store(5)
	src[0][0].TopIPs.Increment("1.1.1.1", 3)
	src[0][0].TopCountries.Increment("RO", 2)

	src[5][12].RecordsPerPeriod.Store(7)
	src[5][12].TopIPs.Increment("8.8.8.8", 4)
	src[5][12].TopBrowsers.Increment(Chrome, 1)

	src[30][23].RecordsPerPeriod.Store(1)
	src[30][23].TopASN.Increment("AS1234", 9)

	var dst MonthActive

	src.deepCopyInto(&dst)

	// --- Verify copied values match source ---
	require.Equal(t, uint32(5), dst[0][0].RecordsPerPeriod.Load())
	require.Equal(t, uint32(3), dst[0][0].TopIPs.Count("1.1.1.1"))
	require.Equal(t, uint32(2), dst[0][0].TopCountries.Count("RO"))

	require.Equal(t, uint32(7), dst[5][12].RecordsPerPeriod.Load())
	require.Equal(t, uint32(4), dst[5][12].TopIPs.Count("8.8.8.8"))
	require.Equal(t, uint32(1), dst[5][12].TopBrowsers.Count(Chrome))

	require.Equal(t, uint32(1), dst[30][23].RecordsPerPeriod.Load())
	require.Equal(t, uint32(9), dst[30][23].TopASN.Count("AS1234"))

	// --- Verify empty slots remain empty ---
	require.Equal(t, uint32(0), dst[0][0].TopASN.Count("AS1234"))
	require.Equal(t, uint32(0), dst[10][10].RecordsPerPeriod.Load())

	// --- CRITICAL: deep copy isolation ---
	// Modifying src must NOT affect dst
	src[0][0].RecordsPerPeriod.Store(999)
	src[0][0].TopIPs.Increment("1.1.1.1", 10) // add to existing

	require.Equal(t,
		uint32(5),
		dst[0][0].RecordsPerPeriod.Load(),

		"dst RecordsPerPeriod corrupted by src mutation",
	)
	require.Equal(t,
		uint32(3),
		dst[0][0].TopIPs.Count("1.1.1.1"),

		"dst TopIPs corrupted by src mutation",
	)
}
