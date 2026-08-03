package datacenter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/analytics77/domain/analytics"
	"github.com/tudorhulban/analytics77/domain/dhelpers"
	"github.com/tudorhulban/analytics77/shared"
)

func TestAddEvents_OlderThanPreviousMonth_DoesNotPanic(t *testing.T) {
	dc := NewDataCenter()

	req := shared.ParamsAddEvent{
		SiteKey:       "localhost",
		IP:            "127.0.0.1",
		TimestampUNIX: 0, // 1970 — older than any retained window
		IsPrivateIP:   true,
	}

	require.NotPanics(t,
		func() {
			dc.AddEvents(&req)
		},
	)
}

func TestAddEventsAcrossBoundaries(t *testing.T) {
	site := "localhost"

	now := time.Now()

	tests := []struct {
		producedAtUTC time.Time
		description   string

		localOffsetHours  int8
		utcExpectedDay    int8
		utcExpectedHour   int8
		expectedLastMonth bool
	}{
		{
			description:   "1.Previous hour inside same day",
			producedAtUTC: time.Date(now.Year(), now.Month(), now.Day(), 15, 0, 0, 0, time.UTC), // 15:00 UTC

			localOffsetHours: +3,
			utcExpectedDay:   int8(now.Day()),
			utcExpectedHour:  18,
		},
		{
			description:   "2.Previous hour crosses midnight (day → day-1)",
			producedAtUTC: time.Date(now.Year(), now.Month(), now.Day(), 0, 30, 0, 0, time.UTC), // 00:30 UTC

			localOffsetHours: -1,
			utcExpectedDay:   int8(now.Day()) - 1,
			utcExpectedHour:  23,
		},
		{
			description:   "3.Previous hour across month boundary (day=1 → day=31)",
			producedAtUTC: time.Date(now.Year(), now.Month(), 1, 0, 15, 0, 0, time.UTC),

			localOffsetHours:  -1,
			utcExpectedDay:    int8(time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, time.UTC).Day()), //nolint:revive
			utcExpectedHour:   23,
			expectedLastMonth: true,
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.description,
			func(t *testing.T) {
				dc := NewDataCenter()
				require.NotNil(t, dc)

				req := shared.ParamsAddEvent{
					SiteKey:         site,
					IP:              "127.0.0.1",
					Browser:         analytics.Brave,
					OperatingSystem: analytics.Linux,

					TimestampUNIX:  tc.producedAtUTC.Unix(),
					OffsetUTCHours: int64(tc.localOffsetHours),

					IsPrivateIP: true,
				}

				require.Empty(t, dc.AddEvents(&req))

				registry := dc.GetRegistry(Site(site))
				require.NotNil(t, registry)
				require.NotZero(t, registry.CalendarMonthCurrentNumber.Load())

				var (
					days   [31]int8
					noDays int8
				)

				if !tc.expectedLastMonth {
					days, noDays = registry.CurrentMonthDaysWithData()
					require.EqualValues(t,
						1,
						noDays,

						"number of days not matching, got:%d",
						noDays,
					)
				} else {
					days, noDays = registry.PreviousMonthDaysWithData()
					require.EqualValues(t,
						1,
						noDays,

						"number of days not matching, got:%d",
						noDays,
					)
				}

				require.EqualValues(t,
					tc.utcExpectedDay,
					dhelpers.IndexToCalendarDay(days[0]),

					"day not matching, expected:%d got:%d",
					tc.utcExpectedDay,
					dhelpers.IndexToCalendarDay(days[0]),
				)

				var (
					hours   [24]int8
					noHours int8
				)

				if !tc.expectedLastMonth {
					hours, noHours = registry.CurrentDayHoursWithData(tc.producedAtUTC.Unix(), int64(tc.localOffsetHours))
					require.EqualValues(t, noHours, 1)
					assert.EqualValues(t,
						tc.utcExpectedHour,
						hours[0],

						"hours not matching",
					)
				} else {
					hours, noHours = registry.PreviousMonthHoursWithData(tc.utcExpectedDay)
					require.EqualValues(t,
						noHours,
						1,

						"for day:%d",
						days[0],
					)
					assert.EqualValues(t,
						tc.utcExpectedHour,
						hours[0],

						"hours not matching, expected:%d got:%d",
						tc.utcExpectedHour,
						hours[0],
					)
				}
			},
		)
	}
}
