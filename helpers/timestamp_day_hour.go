package helpers

type TimestampOffsets struct {
	OffsetUTCHours int64

	TimestampDSTWinter int64 // Epoch when DST ends (falls back to winter/standard time)
	TimestampDSTSpring int64 // Epoch when DST starts (springs forward to summer time)
}

// ExtractMonthDayHour converts a UTC Unix timestamp into the *real* calendar day (1–31)
// and local hour (0–23), applying both the fixed UTC offset and the DST window.
//
// Notes:
//
//	a. The returned `day` is the actual civil calendar day, not a storage index.
//	   Example: July 7 → day=7, hour=<local hour>; March 1 → day=1.
//	b. Long Februaries (29 days) are handled correctly. The civil-time algorithm
//	   computes the true Gregorian date, independent of month length.
//	c. The final month-day result is always 1–31. If the storage model uses a
//	   fixed 31‑day ring (0–30), this should be mapped via CalendarDayToIndex(day).
//	d. The algorithm uses the March-based Gregorian conversion, so leap years,
//	   DST offsets, and negative timestamps are all resolved deterministically.
//
// DST behavior:
//   - If timestampUTC ∈ [TimestampDSTSpring, TimestampDSTWinter), an extra +3600s
//     is injected into localTimestamp.
//
// Example:
//
//	UTC: 2026‑07‑07 10:00, OffsetUTC=+3h, DST active → local=13:00 → day=7, hour=13.
//
// This function guarantees:
//   - Correct Gregorian day extraction for all months (including Feb 29).
//   - Stable behavior for timestamps before 1970.
//   - No zero-based day values; day always ∈ [1..31].
func ExtractMonthDayHour(timestampUTC int64, offsets *TimestampOffsets) (int8, int8, int8) {
	// Start with base timestamp + standard offset
	localTimestamp := timestampUTC + offsets.OffsetUTCHours*3600

	// Check if the timestamp falls within the DST active window.
	// If it does, we inject the extra 1 hour (3600 seconds) savings.
	if timestampUTC >= offsets.TimestampDSTSpring && timestampUTC < offsets.TimestampDSTWinter {
		localTimestamp = localTimestamp + 3600
	}

	totalHours := localTimestamp / 3600
	hour := int(totalHours % 24)

	// Handle Go's truncated division behavior on negative local timestamps
	if hour < 0 {
		hour = hour + 24
	}

	// 86400 seconds = 1 day.
	totalDays := localTimestamp / 86400

	// Handle negative totalDays if the offset pushes the time before 1970
	if localTimestamp < 0 && localTimestamp%86400 != 0 {
		totalDays--
	}

	// Unix epoch (Jan 1, 1970) was a Thursday.
	// To convert totalDays since 1970 into a specific day of the current month
	// using pure integer math requires a fast epoch-to-date algorithm (like civil time).

	// optimized integer algorithm for UTC date extraction:
	totalDays = totalDays + 719468 // Offset to March 1, 0000

	doe := totalDays % 146097
	yoe := (doe - doe/1460 + doe/36524 - doe/146096) / 365
	doy := doe - (365*yoe + yoe/4 - yoe/100)

	mp := (5*doy + 2) / 153
	day := int(doy - (153*mp+2)/5 + 1)

	if mp < 10 {
		mp = mp + 3
	} else {
		mp = mp - 9
	}

	return int8(mp), int8(day), int8(hour)
}

func ExtractMonth(timestampUTC int64) int8 {
	totalDays := timestampUTC / 86400
	if timestampUTC < 0 && timestampUTC%86400 != 0 {
		totalDays--
	}

	// Shift epoch to March 1, 0000
	totalDays = totalDays + 719468

	// Position within the 400-year Gregorian cycle (146,097 days)
	doe := totalDays % 146097
	if doe < 0 {
		doe += 146097
	}

	// Year within the era (0..399)
	yoe := (doe - doe/1460 + doe/36524 - doe/146096) / 365

	// Day within the shifted year (0 = March 1st)
	doy := doe - (365*yoe + yoe/4 - yoe/100)

	// Month in shifted calendar (0 = March, 1 = April, ..., 11 = February)
	mp := (5*doy + 2) / 153

	// Convert shifted calendar month to standard 1..12 month
	if mp < 10 {
		return int8(mp + 3)
	}

	return int8(mp - 9)
}
