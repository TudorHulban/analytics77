package datacenter

import (
	"fmt"
	"slices"
	"strings"
	"time"

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

func (dc *DataCenter) GetPreviousHourRecordsPerSite(offsets *helpers.TimestampOffsets) ResponseRecordsPerSite {
	_, ixDay, ixHour := helpers.ExtractMonthDayHour(
		time.Now().Unix(),
		offsets,
	)

	prevHour := ixHour - 1
	day := ixDay

	if prevHour < 0 {
		prevHour = 23

		day = ixDay - 1
		if day < 1 {
			day = 31 // last day of the 31-day model
		}
	}

	dc.mu.Lock()

	result := make(map[string]uint32, len(dc.data))

	for siteKey, registry := range dc.data {
		result[string(siteKey)] = registry.
			GetCurrentMonth()[dhelpers.CalendarDayToIndex(day)][prevHour].RecordsPerPeriod.Load()
	}

	dc.mu.Unlock()

	return result
}

func (dc *DataCenter) GetCurrentHourRecordsPerSite(offsets *helpers.TimestampOffsets) ResponseRecordsPerSite {
	_, ixDay, ixHour := helpers.ExtractMonthDayHour(
		time.Now().Unix(),
		offsets,
	)

	dc.mu.Lock()

	result := make(map[string]uint32, len(dc.data))

	for siteKey, registry := range dc.data {
		result[string(siteKey)] = registry.
			GetCurrentMonth()[dhelpers.CalendarDayToIndex(ixDay)][ixHour].RecordsPerPeriod.Load()
	}

	dc.mu.Unlock()

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
