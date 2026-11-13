package datacenter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/analytics77/domain/analytics"
	"github.com/tudorhulban/analytics77/domain/offsets"
	"github.com/tudorhulban/analytics77/domain/site"
	"github.com/tudorhulban/analytics77/fixtures"
	"github.com/tudorhulban/analytics77/shared"
	"github.com/tudorhulban/hxhelpers"
)

func TestPreviousHourRecords(t *testing.T) {
	dc, errCr := NewDataCenter(fixtures.DefaultOffset())
	require.NoError(t, errCr)
	require.NotNil(t, dc)

	siteName := "localhost"

	localTime := time.Now()
	_, offsetSeconds := localTime.Zone()

	utcOffsetHours := offsetSeconds / 3600

	hourLocalAsUTC := offsets.GetHourInOffsetAsUTC(
		localTime,
		localTime.Hour(),
		utcOffsetHours,
	)

	t.Logf(
		"hourLocalAsUTC: %s",
		time.Unix(hourLocalAsUTC, 0).UTC().Format("2006-01-02 15:04:05"),
	)

	req1 := shared.ParamsAddEvent{
		SiteKey: siteName,

		IP:              "127.0.0.1",
		Browser:         analytics.Brave,
		OperatingSystem: analytics.Linux,

		TimestampUNIX: hourLocalAsUTC,
		IsPrivateIP:   true,
	}

	require.Empty(t, dc.AddEvents(&req1))

	registryBefore, exists := dc.GetRegistry(site.SiteName(siteName))
	require.True(t, exists)
	require.EqualValues(t,
		localTime.Month(),
		registryBefore.CalendarMonthCurrentNumber.Load(),
	)

	recordsPreviousHourBefore := dc.GetPreviousHourRecordsPerSite(
		&offsets.TimestampOffsets{
			OffsetUTCHours: int64(utcOffsetHours),
		},
	)
	require.Len(t, recordsPreviousHourBefore, 1)
	require.Zero(t,
		recordsPreviousHourBefore[site.SiteName(siteName)],

		"previous hour has no data yet",
	)

	req2 := shared.ParamsAddEvent{
		SiteKey:         siteName,
		IP:              "127.0.0.1",
		Browser:         analytics.Brave,
		OperatingSystem: analytics.Linux,

		TimestampUNIX: hourLocalAsUTC - 1*3600,
		IsPrivateIP:   true,
	}

	require.Empty(t, dc.AddEvents(&req2))

	registryAfter, exists := dc.GetRegistry(site.SiteName(siteName))
	require.True(t, exists)

	hours, howMany := registryAfter.CurrentDayHoursWithData(localTime.UTC())
	require.EqualValues(t, 2, howMany)
	require.Empty(t,
		hxhelpers.NotInSliceSource(
			hours[:],
			int8(localTime.Hour()),
			int8(localTime.Hour())),
	)

	recordsPreviousHourAfter := dc.GetPreviousHourRecordsPerSite(
		&offsets.TimestampOffsets{
			OffsetUTCHours: int64(utcOffsetHours),
		},
	)
	require.Len(t, recordsPreviousHourAfter, 1)
	require.EqualValues(t,
		1,
		recordsPreviousHourAfter[site.SiteName(siteName)],

		"previous hour should have req2 now",
	)
}
