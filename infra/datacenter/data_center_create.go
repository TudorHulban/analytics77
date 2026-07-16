package datacenter

import (
	"github.com/tudorhulban/analytics77/domain/analytics"
	"github.com/tudorhulban/analytics77/domain/dhelpers"
	"github.com/tudorhulban/analytics77/helpers"
	"github.com/tudorhulban/analytics77/shared"
)

// AddEvents ingests one or more events into the DataCenter.
//
// Pipeline semantics:
//
//  1. Validation (no locks)
//     - Each event is validated independently.
//     - Invalid events are collected into errorsBatch.
//     - Valid event indexes are stored for processing.
//
//  2. Pointer resolution under lock (structure only)
//     - Ensures dc.data[site] exists.
//     - Extracts (calendarMonth, calendarDay, hour) using site-level DST/offsets.
//     - Determines the correct month buffer:
//     * MonthCurrent   → events from the active calendar month
//     * MonthPrevious  → events from the immediately preceding month
//     * Rollover()     → triggered when an event arrives from the next month
//     - Writes resolved[outIx] as a direct pointer to the correct [day][hour] slot.
//     - No increments or allocations occur under the lock.
//
//  3. Fast-path increments (outside lock)
//     - Atomic increments into the resolved MetricActive slots.
//     - Updates TopIPs, TopBrowsers, TopASN, TopCountries, TopCities.
//
//  4. Error return
//     - If any event failed validation, returns the collected errors.
//     - Otherwise returns nil.
//
// Month routing rules:
//
//	evMonth == currentMonth
//	    → MonthCurrent
//
//	evMonth + 1 == currentMonth   (with December → January wrap)
//	    → MonthPrevious
//
//	evMonth - 1 == currentMonth
//	    → Rollover() then MonthCurrent
//
//	otherwise
//	    → event is older than previous month and ignored.
//
// No time mocking, no abstraction layers, no background rollover.
// Month transitions are strictly event-driven.
func (dc *DataCenter) AddEvents(events ...*shared.ParamsAddEvent) []error {
	if len(events) == 0 {
		return nil
	}

	errorsBatch := make([]error, 0)
	indexesNoError := make([]int, 0, len(events))

	var hasErrors bool

	// Validate without locks
	for ix, event := range events {
		if event == nil {
			continue
		}

		if errorsValidation := event.Validate(); errorsValidation != nil {
			hasErrors = true

			errorsBatch = append(
				errorsBatch,
				errorsValidation...)

			continue
		}

		indexesNoError = append(indexesNoError, ix)
	}

	// Resolve registry pointers under lock (structure only)
	//
	//    We do NOT increment inside the lock.
	//    We only ensure that dc.data[site] exists and is stable.
	//
	resolved := make([]*analytics.MetricActive, len(indexesNoError))

	dc.mu.Lock()
	for outIx, eventIndex := range indexesNoError {
		event := events[eventIndex]
		site := Site(event.SiteKey)

		registrySite, exists := dc.data[site]
		if !exists {
			registrySite = analytics.NewRegistry()
			dc.data[site] = registrySite
		}

		evCalendarMonth, evCalendarDay, evHour := helpers.ExtractMonthDayHour(
			event.TimestampUNIX,
			&helpers.TimestampOffsets{
				OffsetUTCHours: event.OffsetUTCHours,

				TimestampDSTSpring: registrySite.TimestampDSTSpring,
				TimestampDSTWinter: registrySite.TimestampDSTWinter,
			},
		)

		if registrySite.CalendarMonthCurrentNumber == 0 {
			registrySite.CalendarMonthCurrentNumber = evCalendarMonth
		}

		switch {
		case evCalendarMonth == registrySite.CalendarMonthCurrentNumber:
			// Same calendar month → current month buffer.
			resolved[outIx] = &registrySite.
				GetCurrentMonth()[dhelpers.CalendarDayToIndex(evCalendarDay)][evHour]

		case evCalendarMonth+1 == registrySite.CalendarMonthCurrentNumber || (evCalendarMonth == 12 && registrySite.CalendarMonthCurrentNumber == 1):
			// Event belongs to previous calendar month.
			resolved[outIx] = &registrySite.
				GetPreviousMonth()[dhelpers.CalendarDayToIndex(evCalendarDay)][evHour]

		case evCalendarMonth-1 == registrySite.CalendarMonthCurrentNumber:
			// Event belongs to next calendar month → rollover required.
			registrySite.Rollover()

			// After rollover, the new month becomes current.
			resolved[outIx] = &registrySite.
				GetCurrentMonth()[dhelpers.CalendarDayToIndex(evCalendarDay)][evHour]

		default:
			// Older than previous month → either archive or ignore.
			// For now, ignore in fast path.
			continue
		}
	}
	dc.mu.Unlock()

	// Perform increments OUTSIDE the lock
	//
	//    This is the fast path:
	//    - atomic increments
	//    - no allocations
	//    - no global lock
	//
	for outIx, eventIndex := range indexesNoError {
		slot := resolved[outIx]
		if slot == nil {
			continue // event already ignored in default above
		}

		event := events[eventIndex]

		slot.RecordsPerPeriod.Add(1)

		slot.TopIPs.Increment(event.IP, 1)
		slot.TopBrowsers.Increment(event.Browser, 1)
		slot.TopASN.Increment(event.ASNOrganization, 1)
		slot.TopCountries.Increment(event.Country, 1)
		slot.TopCities.Increment(event.City, 1)
	}

	if hasErrors {
		return errorsBatch
	}

	return nil
}
