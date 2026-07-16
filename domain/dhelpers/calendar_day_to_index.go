package dhelpers

// CalendarDayToIndex converts 1–31 → 0–30
func CalendarDayToIndex(dayCalendar int8) int8 {
	return dayCalendar - 1
}

// IndexToCalendarDay converts 0–30 → 1–31
func IndexToCalendarDay(dayIndex int8) int8 {
	return dayIndex + 1
}
