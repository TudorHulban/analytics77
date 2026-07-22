package analytics

import (
	"github.com/tudorhulban/analytics77/domain/dhelpers"
	"github.com/tudorhulban/analytics77/helpers"
)

// PreviousMonthDaysWithData returns a stack‑allocated array of up to 31 days.
// The second return value is the number of valid entries in the array.
//
// This avoids heap allocation: callers must iterate only up to `count`.
func (r *Registry) PreviousMonthDaysWithData() ([31]int8, int8) {
	var (
		out   [31]int8
		count int8
	)

	previousSlot := r.GetPreviousSlot()

	for day := range int8(31) {
		// Scan hours for any non‑zero record
		for hour := range int8(24) {
			m := (*previousSlot).GetMetric(day, hour)

			if (*m).GetRecordsPerPeriod() != 0 {
				out[count] = day
				count++

				break
			}
		}
	}

	return out, count
}

// NOTE ABOUT DAY INDEXING
//
// All presence / lookup methods in this Registry work with a **zero-based day index**.
// This matches the internal storage layout:
//
//     MonthCurrent[0]  → Day 1
//     MonthCurrent[1]  → Day 2
//     ...
//     MonthCurrent[30] → Day 31
//
// Ingestion APIs use **1-based calendar days** (DayOfMonth = 1..31),
// but internal arrays use **0-based indexing** for performance and simplicity.
//
// Method must convert:
//
//     dayIndex = DayOfMonth - 1
//
// Failing to apply this offset will return empty results even when data exists.
// This is a common source of subtle bugs in tests and callers.

// PreviousMonthHoursWithData returns a stack‑allocated array of up to 24 hours
// that contain at least one non‑zero metric for the given day.
// The second return value is the number of valid entries.
//
// This avoids heap allocation: callers must iterate only up to `count`.
func (r *Registry) PreviousMonthHoursWithData(forCalendarDay int8) ([24]int8, int8) {
	var (
		hoursWithData [24]int8
		count         int8
	)

	if forCalendarDay < 1 || forCalendarDay > 31 {
		return hoursWithData, 0
	}

	dayIndex := dhelpers.CalendarDayToIndex(forCalendarDay)
	prev := r.GetPreviousSlot()

	for hour := int8(0); hour < 24; hour++ {
		slot := (*prev).GetMetric(dayIndex, hour)

		if (*slot).GetRecordsPerPeriod() != 0 {
			hoursWithData[count] = hour
			count++
		}
	}

	return hoursWithData, count
}

// CurrentMonthHoursWithData returns a stack‑allocated array of up to 24 hours
// that contain at least one non‑zero metric for the given day.
// The second return value is the number of valid entries.
//
// This avoids heap allocation: callers must iterate only up to `count`.
func (r *Registry) CurrentMonthHoursWithData(forCalendarDay int8) ([24]int8, int8) {
	var (
		hoursWithData [24]int8
		count         int8
	)

	if forCalendarDay < 1 || forCalendarDay > 31 {
		return hoursWithData, 0
	}

	dayIndex := dhelpers.CalendarDayToIndex(forCalendarDay)
	curr := r.GetActiveSlot()

	for hour := int8(0); hour < 24; hour++ {
		slot := (*curr).GetMetric(dayIndex, hour)

		if (*slot).GetRecordsPerPeriod() != 0 {
			hoursWithData[count] = hour
			count++
		}
	}

	return hoursWithData, count
}

// Registry stores storage days so
// the day numbers need to be translated to calendar days.
//
// Registry stores storage hours as UTC time.
func (r *Registry) CurrentDayHoursWithData(timestampUTC, offsetUTCHours int64) ([24]int8, int8) {
	_, currentDay, _ := helpers.ExtractMonthDayHour(
		timestampUTC,
		&helpers.TimestampOffsets{
			OffsetUTCHours: offsetUTCHours,

			TimestampDSTSpring: r.TimestampDSTSpring,
			TimestampDSTWinter: r.TimestampDSTWinter,
		},
	)

	return r.CurrentMonthHoursWithData(currentDay)
}

// CurrentMonthDaysWithData returns a stack‑allocated array of up to 31 days
// that contain at least one non‑zero metric. The second return value is the
// number of valid entries.
//
// This avoids heap allocation: callers must iterate only up to `count`.
//
// Registry stores storage days so
// the day numbers need to be translated to calendar days.
func (r *Registry) CurrentMonthDaysWithData() ([31]int8, int8) {
	var (
		daysWithData [31]int8
		howMany      int8
	)

	activeSlot := r.GetActiveSlot()

	for day := range int8(31) {
		for hour := range int8(24) {
			slot := (*activeSlot).GetMetric(day, hour)

			if (*slot).GetRecordsPerPeriod() != 0 {
				daysWithData[howMany] = day
				howMany++

				break
			}
		}
	}

	return daysWithData, howMany
}
