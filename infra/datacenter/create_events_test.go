package datacenter

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/analytics77/domain/analytics"
	"github.com/tudorhulban/analytics77/domain/calendar"
	"github.com/tudorhulban/analytics77/domain/site"
	"github.com/tudorhulban/analytics77/fixtures"
	"github.com/tudorhulban/analytics77/shared"
)

func TestAddEvents_HappyPath(t *testing.T) {
	now := time.Now()

	siteName := "example.com"

	dc, errCr := NewDataCenter(
		fixtures.DefaultOffset(now),
	)
	require.NoError(t, errCr)
	require.NotNil(t, dc)

	req := shared.ParamsAddEvent{
		SiteKey: siteName,
		Country: "ROU",
		City:    "Iasi",

		IP:            "82.77.237.37",
		TimestampUNIX: now.Unix(),
	}

	require.Empty(t, dc.AddEvents(&req))
	require.Len(t, dc.data, 1)
	require.EqualValues(t, dc.GetSiteNames()[0], siteName)

	dcSite, exists := dc.data[site.SiteName(siteName)]
	require.True(t, exists)
	require.NotNil(t, dcSite)
	require.NotNil(t, dcSite.Data)

	dcSite.Data.Snapshot(os.Stdout)
}

func TestAddEvents_OlderThanPreviousMonth_DoesNotPanic(t *testing.T) {
	dc, errCr := NewDataCenter(fixtures.DefaultOffset())
	require.NoError(t, errCr)
	require.NotNil(t, dc)

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

func TestAddEvents_AcrossBoundaries(t *testing.T) {
	siteTest := site.Site{
		Name:             site.SiteName(t.Name()),
		TimestampOffsets: fixtures.DefaultOffset(),
	}

	now := time.Now()

	tests := []struct {
		producedAtUTC time.Time
		description   string

		utcExpectedDay    int8
		utcExpectedHour   int8
		expectedLastMonth bool
	}{
		// {
		// 	description: "1.Previous hour inside same day",
		// 	producedAtUTC: time.Date(
		// 		now.Year(),
		// 		now.Month(),
		// 		now.Day(),
		// 		now.Hour()-1, // UTC previous hour
		// 		0, 0, 0,
		// 		time.UTC,
		// 	),

		// 	utcExpectedDay:  int8(now.Day()),
		// 	utcExpectedHour: int8(now.Hour()) - 1,
		// },
		// {
		// 	description: "2.Midnight hour",
		// 	producedAtUTC: time.Date(
		// 		now.Year(),
		// 		now.Month(),
		// 		now.Day(),
		// 		0,
		// 		30, 0, 0, // 00:30 UTC
		// 		time.UTC,
		// 	),

		// 	utcExpectedDay:  int8(now.Day()),
		// 	utcExpectedHour: 0,
		// },
		// {
		// 	description: "3.Previous hour before midnight (day → day-1)",
		// 	producedAtUTC: time.Date(
		// 		now.Year(),
		// 		now.Month(),
		// 		now.Day()-1,
		// 		23,
		// 		30, 0, 0, // 23:30 UTC
		// 		time.UTC,
		// 	),

		// 	utcExpectedDay:  int8(now.Day()) - 1,
		// 	utcExpectedHour: 23,
		// },
		{
			description: "4.Previous hour before month boundary (day=1 → day=31)",
			producedAtUTC: time.Date(
				now.Year(),
				now.Month(),
				1,
				0,
				30, 0, 0, // 00:30 UTC
				time.UTC,
			),

			utcExpectedDay: int8(time.Date(
				now.Year(),
				now.Month(),
				1,
				0,
				0, 0, 0,
				time.UTC).Day(),
			), //nolint:revive

			utcExpectedHour:   23,
			expectedLastMonth: true,
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.description,

			func(t *testing.T) {
				datacenter, errCr := NewDataCenter(
					fixtures.DefaultOffset(now),
				)
				require.NoError(t, errCr)
				require.NotNil(t, datacenter)

				req := shared.ParamsAddEvent{
					SiteKey:         string(siteTest.Name),
					IP:              "127.0.0.1",
					Browser:         analytics.Brave,
					OperatingSystem: analytics.Linux,

					TimestampUNIX: tc.producedAtUTC.Unix(),
					IsPrivateIP:   true,
				}

				require.Empty(t, datacenter.AddEvents(&req))

				registry, exists := datacenter.GetRegistry(siteTest.Name)
				require.True(t, exists)
				require.NotNil(t, registry)
				require.NotZero(t, registry.CalendarMonthCurrentNumber.Load())

				fmt.Println(datacenter.Snapshot(os.Stdout))

				var (
					days   [31]int8
					noDays int8
				)

				if tc.expectedLastMonth {
					days, noDays = registry.PreviousMonthDaysWithData()
					require.EqualValues(t,
						1,
						noDays,

						"number of days not matching, got:%d",
						noDays,
					)
				} else {
					days, noDays = registry.CurrentMonthDaysWithData()
					require.EqualValues(t,
						1,
						noDays,

						"number of days not matching, got:%d",
						noDays,
					)
				}

				assert.EqualValues(t,
					tc.utcExpectedDay,
					calendar.IndexToCalendarDay(days[0]),

					"day not matching, expected:%d got:%d",
					tc.utcExpectedDay,
					calendar.IndexToCalendarDay(days[0]),
				)

				var (
					hours   [24]int8
					noHours int8
				)

				if !tc.expectedLastMonth {
					hours, noHours = registry.
						CurrentDayHoursWithData(
							tc.producedAtUTC,
						)
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
