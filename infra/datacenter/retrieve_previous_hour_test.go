package datacenter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/analytics77/domain/analytics"
	"github.com/tudorhulban/analytics77/helpers"
	"github.com/tudorhulban/analytics77/shared"
	"github.com/tudorhulban/hxhelpers"
)

func TestPreviousHourRecords(t *testing.T) {
	dc := NewDataCenter()
	require.NotNil(t, dc)

	site := "localhost"

	localTime := time.Now()
	_, offsetSeconds := localTime.Zone()

	utcOffsetHours := offsetSeconds / 3600

	hourLocalAsUTC := helpers.GetHourInOffsetAsUTC(
		localTime,
		localTime.Hour(),
		utcOffsetHours,
	)

	t.Logf(
		"hourLocalAsUTC: %s",
		time.Unix(hourLocalAsUTC, 0).UTC().Format("2006-01-02 15:04:05"),
	)

	req1 := shared.ParamsAddEvent{
		SiteKey: site,

		IP:              "127.0.0.1",
		Browser:         analytics.Brave,
		OperatingSystem: analytics.Linux,

		TimestampUNIX:  hourLocalAsUTC,
		OffsetUTCHours: int64(utcOffsetHours),

		IsPrivateIP: true,
	}

	require.Empty(t, dc.AddEvents(&req1))

	registryBefore := dc.GetRegistry(Site(site))
	require.EqualValues(t,
		localTime.Month(),
		registryBefore.CalendarMonthCurrentNumber.Load(),
	)

	recordsPreviousHourBefore := dc.GetPreviousHourRecordsPerSite(
		&helpers.TimestampOffsets{
			OffsetUTCHours: int64(utcOffsetHours),
		},
	)
	require.Len(t, recordsPreviousHourBefore, 1)
	require.Zero(t,
		recordsPreviousHourBefore[site],

		"previous hour has no data yet",
	)

	req2 := shared.ParamsAddEvent{
		SiteKey:         site,
		IP:              "127.0.0.1",
		Browser:         analytics.Brave,
		OperatingSystem: analytics.Linux,

		TimestampUNIX:  hourLocalAsUTC - 1*3600,
		OffsetUTCHours: int64(utcOffsetHours),

		IsPrivateIP: true,
	}

	require.Empty(t, dc.AddEvents(&req2))

	registryAfter := dc.GetRegistry(Site(site))
	hours, howMany := registryAfter.CurrentDayHoursWithData(localTime.Unix(), int64(utcOffsetHours))
	require.EqualValues(t, 2, howMany)
	require.Empty(t,
		hxhelpers.NotInSliceSource(
			hours[:],
			int8(localTime.Hour()),
			int8(localTime.Hour())),
	)

	recordsPreviousHourAfter := dc.GetPreviousHourRecordsPerSite(
		&helpers.TimestampOffsets{
			OffsetUTCHours: int64(utcOffsetHours),
		},
	)
	require.Len(t, recordsPreviousHourAfter, 1)
	require.EqualValues(t,
		1,
		recordsPreviousHourAfter[site],

		"previous hour should have req2 now",
	)
}
