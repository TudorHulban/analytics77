package datacenter

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/analytics77/domain/analytics"
	"github.com/tudorhulban/analytics77/domain/calendar"
	"github.com/tudorhulban/analytics77/domain/site"
)

type ResponseRecordsPerSite map[site.SiteName]uint32

func (r ResponseRecordsPerSite) Verify(tb testing.TB, siteNames []site.SiteName, expectedNoRecords []uint32) {
	tb.Helper()

	require.Equal(tb,
		len(siteNames),
		len(expectedNoRecords),
	)

	for ix, key := range siteNames {
		value, exists := r[key]
		require.True(tb, exists)
		require.Equal(tb,
			value,
			expectedNoRecords[ix],

			"value (%v) vs. expected (%v)",
			value,
			expectedNoRecords[ix],
		)
	}
}

func (r ResponseRecordsPerSite) String() string {
	if len(r) == 0 {
		return "{}"
	}

	keys := make([]site.SiteName, 0, len(r))

	for k := range r {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	var builder strings.Builder
	builder.WriteString("{")

	for ix, key := range keys {
		fmt.Fprintf(&builder, "%s: %d",

			key,
			r[key],
		)

		if ix < len(keys)-1 {
			builder.WriteString(", ")
		}
	}

	builder.WriteString("}")

	return builder.String()
}

func (dc *DataCenter) PreviousHourRecordsPerSiteAt(nowUTC time.Time) ResponseRecordsPerSite {
	previousHourUTC := nowUTC.Add(-1 * time.Hour)

	dc.mu.RLock()

	result := make(map[site.SiteName]uint32, len(dc.data))

	for siteName, dcSite := range dc.data {
		prevMonth := int8(previousHourUTC.Month()) //nolint:gosec
		prevDay := int8(previousHourUTC.Day())     //nolint:gosec
		prevHour := previousHourUTC.Hour()         // TODO: move definition outside loop for lower allocs

		var month *analytics.MonthActive

		switch prevMonth {
		case int8(dcSite.Data.CalendarMonthCurrentNumber.Load()): //nolint:gosec
			month = dcSite.Data.GetActiveSlot()

		default:
			// The hour we're after fell into the previous calendar month.
			month = dcSite.Data.GetPreviousSlot()
		}

		result[siteName] = month[calendar.CalendarDayToIndex(prevDay)][prevHour].RecordsPerPeriod.Load()
	}

	dc.mu.RUnlock()

	return result
}

// GetPreviousHourRecordsPerSite returns, per site, the record count for the
// hour immediately preceding the current one.
func (dc *DataCenter) GetPreviousHourRecordsPerSite() ResponseRecordsPerSite {
	return dc.
		PreviousHourRecordsPerSiteAt(
			time.Now(),
		)
}

func (dc *DataCenter) GetCurrentHourRecordsPerSite() ResponseRecordsPerSite {
	nowUTC := time.Now().UTC()

	ixDay := int8(nowUTC.Day()) //nolint:gosec
	ixHour := nowUTC.Hour()

	dc.mu.RLock()

	result := make(map[site.SiteName]uint32, len(dc.data))

	for siteName, dcSite := range dc.data {
		result[siteName] = dcSite.
			Data.
			GetActiveSlot()[calendar.CalendarDayToIndex(ixDay)][ixHour].
			RecordsPerPeriod.Load()
	}

	dc.mu.RUnlock()

	return result
}

// CurrentDayHoursWithData returns UTC hours.
func (dc *DataCenter) CurrentDayHoursWithData() ([24]int8, int8) {
	hoursMap := make(map[int8]struct{})

	sites := dc.GetSiteNames()

	// Grab the local/machine Unix timestamp instead of forcing UTC
	now := time.Now()

	for _, siteName := range sites {
		registry, _ := dc.GetRegistry(siteName)

		// Pass the offsets pointer down into the registry method
		hoursWithData, howMany := registry.
			CurrentDayHoursWithData(now.UTC())

		for _, hour := range hoursWithData[:howMany] {
			if _, exists := hoursMap[hour]; !exists {
				hoursMap[hour] = struct{}{}
			}
		}
	}

	var (
		hoursWithData [24]int8
		howMany       int8
	)

	for hour := range hoursMap {
		hoursWithData[howMany] = hour
		howMany++
	}

	return hoursWithData, howMany
}
