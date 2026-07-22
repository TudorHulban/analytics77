package datacenter

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/tudorhulban/analytics77/domain/analytics"
	"github.com/tudorhulban/analytics77/domain/dhelpers"
	"github.com/tudorhulban/analytics77/helpers"
)

type ResponseRecordsPerSite map[string]uint32

func (r ResponseRecordsPerSite) Verify(keys []string, values []uint32) bool {
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

	keys := make([]string, 0, len(r))

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

// GetPreviousHourRecordsPerSite returns, per site, the record count for the
// hour immediately preceding the current one.
func (dc *DataCenter) GetPreviousHourRecordsPerSite(offsets *helpers.TimestampOffsets) ResponseRecordsPerSite {
	return dc.previousHourRecordsPerSiteAt(time.Now().Unix(), offsets)
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
func (dc *DataCenter) previousHourRecordsPerSiteAt(nowUTC int64, offsets *helpers.TimestampOffsets) ResponseRecordsPerSite {
	previousHourUTC := nowUTC - 3600

	dc.mu.RLock()

	result := make(map[string]uint32, len(dc.data))

	for siteKey, registry := range dc.data {
		prevMonth, prevDay, prevHour := helpers.ExtractMonthDayHour(
			previousHourUTC,
			&helpers.TimestampOffsets{
				OffsetUTCHours: offsets.OffsetUTCHours,

				TimestampDSTSpring: registry.TimestampDSTSpring,
				TimestampDSTWinter: registry.TimestampDSTWinter,
			},
		)

		var month *analytics.MonthActive

		switch prevMonth {
		case registry.CalendarMonthCurrentNumber:
			month = registry.GetActiveSlot()

		default:
			// The hour we're after fell into the previous calendar month.
			month = registry.GetPreviousSlot()
		}

		result[string(siteKey)] = month[dhelpers.CalendarDayToIndex(prevDay)][prevHour].RecordsPerPeriod.Load()
	}

	dc.mu.RUnlock()

	return result
}

func (dc *DataCenter) GetCurrentHourRecordsPerSite(offsets *helpers.TimestampOffsets) ResponseRecordsPerSite {
	_, ixDay, ixHour := helpers.ExtractMonthDayHour(
		time.Now().Unix(),
		offsets,
	)

	dc.mu.RLock()

	result := make(map[string]uint32, len(dc.data))

	for siteKey, registry := range dc.data {
		result[string(siteKey)] = registry.GetActiveSlot()[dhelpers.CalendarDayToIndex(ixDay)][ixHour].RecordsPerPeriod.Load()
	}

	dc.mu.RUnlock()

	return result
}

// CurrentDayHoursWithData returns UTC hours.
func (dc *DataCenter) CurrentDayHoursWithData(offsets *helpers.TimestampOffsets) ([24]int8, int8) {
	hoursMap := make(map[int8]struct{})

	sites := dc.GetSiteNames()

	// Grab the local/machine Unix timestamp instead of forcing UTC
	now := time.Now().Unix()

	for _, site := range sites {
		registry := dc.GetRegistry(site)

		// Pass the offsets pointer down into the registry method
		hoursWithData, howMany := registry.CurrentDayHoursWithData(now, offsets.OffsetUTCHours)

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
