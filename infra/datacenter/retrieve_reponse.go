package datacenter

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/tudorhulban/analytics77/domain/analytics"
	"github.com/tudorhulban/analytics77/domain/calendar"
	"github.com/tudorhulban/analytics77/domain/offsets"
	"github.com/tudorhulban/analytics77/domain/site"
)

type ResponseRecordsPerSite map[site.SiteName]uint32

func (r ResponseRecordsPerSite) Verify(keys []site.SiteName, values []uint32) bool {
	if len(keys) != len(values) {
		return false
	}

	for ix, key := range keys {
		value, exists := r[key]

		if !exists || value != values[ix] {
			return false
		}
	}

	return true
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

// previousHourRecordsPerSiteAt does the real work, parameterized on "now" so
// tests can land deterministically on month/DST boundaries.
//
// The previous hour is derived by re-running full month/day/hour extraction
// on (nowUTC - 1h), per site, rather than decrementing the current day/hour
// locally. That's required for two reasons:
//   - month length varies (28-31 days); "day - 1 -> day 31" is only correct
//     for months that actually have 31 days
//   - crossing a DST transition must go through the same civil-calendar math
//     used everywhere else, or the two paths can disagree at the boundary
//
// Each site carries its own DST window, so the extraction uses the caller's
// UTC offset combined with that registry's DST fields — the same convention
// AddEvents uses.
func (dc *DataCenter) previousHourRecordsPerSiteAt(nowUTC int64, offset *offsets.TimestampOffsets) ResponseRecordsPerSite {
	previousHourUTC := nowUTC - 3600

	dc.mu.RLock()

	result := make(map[site.SiteName]uint32, len(dc.data))

	for siteName, dcSite := range dc.data {
		prevMonth, prevDay, prevHour := offsets.ExtractMonthDayHour(
			previousHourUTC,
			&offsets.TimestampOffsets{
				OffsetUTCHours: offset.OffsetUTCHours,

				TimestampDSTSpring: dcSite.Site.TimestampDSTSpring,
				TimestampDSTWinter: dcSite.Site.TimestampDSTWinter,
			},
		)

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
func (dc *DataCenter) GetPreviousHourRecordsPerSite(offset *offsets.TimestampOffsets) ResponseRecordsPerSite {
	return dc.
		previousHourRecordsPerSiteAt(
			time.Now().Unix(),
			offset,
		)
}

func (dc *DataCenter) GetCurrentHourRecordsPerSite(offset *offsets.TimestampOffsets) ResponseRecordsPerSite {
	_, ixDay, ixHour := offsets.ExtractMonthDayHour(
		time.Now().Unix(),
		offset,
	)

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
func (dc *DataCenter) CurrentDayHoursWithData(offset *offsets.TimestampOffsets) ([24]int8, int8) {
	hoursMap := make(map[int8]struct{})

	sites := dc.GetSiteNames()

	// Grab the local/machine Unix timestamp instead of forcing UTC
	now := time.Now().Unix()

	for _, siteName := range sites {
		registry, _ := dc.GetRegistry(siteName)

		// Pass the offsets pointer down into the registry method
		hoursWithData, howMany := registry.
			CurrentDayHoursWithData(
				now,
				offset.OffsetUTCHours,
			)

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
