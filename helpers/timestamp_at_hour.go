package helpers

import "time"

// GetHourInOffsetAsUTC returns the Unix timestamp for a specific hour today at a given UTC offset.
// hour: 0-23 (e.g., 7 for 7 AM, 19 for 7 PM)
// offsetHours: UTC offset in hours (e.g., -5, 1, 3)
func GetHourInOffsetAsUTC(localTime time.Time, hour, utcOffsetHours int) int64 {
	loc := time.FixedZone("CustomZone", utcOffsetHours*3600)

	inZone := localTime.In(loc)

	return time.Date(
		inZone.Year(),
		inZone.Month(),
		inZone.Day(),

		hour, 0, 0, 0, // Minute, second, nanosecond = 0
		loc,
	).
		UTC().Unix()
}
