package datacenter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/analytics77/domain/analytics"
	"github.com/tudorhulban/analytics77/domain/site"
	"github.com/tudorhulban/analytics77/fixtures"
	"github.com/tudorhulban/analytics77/shared"
	"github.com/tudorhulban/hxhelpers"
)

func TestPreviousHourRecords(t *testing.T) {
	datacenter, errCr := NewDataCenter(fixtures.DefaultLocationIANA)
	require.NoError(t, errCr)
	require.NotNil(t, datacenter)

	siteName := "localhost"

	nowUTC := time.Now().UTC()

	t.Logf(
		"hourLocalAsUTC: %s",
		time.Unix(nowUTC.Unix(), 0).UTC().Format("2006-01-02 15:04:05"),
	)

	req1 := shared.ParamsAddEvent{
		SiteKey: siteName,

		IP:              "127.0.0.1",
		Browser:         analytics.Brave,
		OperatingSystem: analytics.Linux,

		TimestampUTC: nowUTC.Unix(),
		IsPrivateIP:  true,
	}

	require.Empty(t, datacenter.AddEvents(&req1))

	registryBefore, exists := datacenter.GetRegistry(site.SiteName(siteName))
	require.True(t, exists)
	require.EqualValues(t,
		nowUTC.Month(),
		registryBefore.CalendarMonthCurrentNumber.Load(),
	)

	recordsPreviousHourBefore := datacenter.GetPreviousHourRecordsPerSite()
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

		TimestampUTC: nowUTC.Unix() - 1*3600,
		IsPrivateIP:  true,
	}

	require.Empty(t, datacenter.AddEvents(&req2))

	registryAfter, exists := datacenter.GetRegistry(site.SiteName(siteName))
	require.True(t, exists)

	hours, howMany := registryAfter.CurrentDayHoursWithData(nowUTC.UTC())
	require.EqualValues(t, 2, howMany)
	require.Empty(t,
		hxhelpers.NotInSliceSource(
			hours[:],
			int8(nowUTC.Hour()),
			int8(nowUTC.Hour())),
	)

	recordsPreviousHourAfter := datacenter.PreviousHourRecordsPerSiteAt(nowUTC.UTC())
	require.Len(t, recordsPreviousHourAfter, 1)
	require.EqualValues(t,
		1,
		recordsPreviousHourAfter[site.SiteName(siteName)],

		"previous hour should have req2 now",
	)
}
