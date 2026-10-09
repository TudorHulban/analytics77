package datacenter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/analytics77/domain/analytics"
	"github.com/tudorhulban/analytics77/shared"
)

func TestDataCenter(t *testing.T) {
	dc := NewDataCenter()
	require.NotNil(t, dc)

	require.Empty(t, dc.AddEvents())
	require.Empty(t, dc.AddEvents(nil))

	site1 := "localhost"

	localTime := time.Now()
	_, offsetSeconds := localTime.Zone()

	utcOffsetHours := offsetSeconds / 3600

	req1 := shared.ParamsAddEvent{
		SiteKey: site1,

		IP:              "127.0.0.1",
		Browser:         analytics.Brave,
		OperatingSystem: analytics.Linux,

		TimestampUNIX:  localTime.Unix(),
		OffsetUTCHours: int64(utcOffsetHours),

		IsPrivateIP: true,
	}

	require.Empty(t, dc.AddEvents(&req1))

	require.Len(t, dc.data, 1, "site1 should have been added")

	registry1, exists1 := dc.data[Site(site1)]
	require.True(t, exists1)
	require.NotNil(t, registry1)

	siteNames := dc.GetSiteNames()
	require.Len(t, siteNames, 1)
	require.EqualValues(t, req1.SiteKey, siteNames[0])

	require.False(t,
		registry1.GetActiveSlot().IsZero(),
	)
	require.True(t,
		registry1.GetPreviousSlot().IsZero(),
	)
	require.True(t, registry1.HistoryAggregateTopN().IsZero())

	hoursWithData, howMany := registry1.CurrentMonthHoursWithData(int8(localTime.Day()))
	require.EqualValues(t, 1, howMany)
	require.EqualValues(t,
		localTime.Hour(),
		hoursWithData[0],
	)

	require.EqualValues(t,
		0,
		registry1.PreviousMonthTotalRecords(),
	)

	require.EqualValues(t,
		1,
		registry1.CurrentMonthTotalRecords(),
	)

	aggregates1 := registry1.CurrentMonthAggregateTopN()
	require.EqualValues(t,
		analytics.Brave,
		*aggregates1.Browsers.Names[0].Load(),

		"browser want:%s got:%s",
		analytics.Brave,
		aggregates1.Browsers.Names[0].Load(),
	)

	dc.Advance(Site(site1))

	require.EqualValues(t,
		0,
		registry1.CurrentMonthTotalRecords(),
	)

	require.EqualValues(t,
		1,
		registry1.PreviousMonthTotalRecords(),
	)

	dc.Advance(Site(site1))

	require.EqualValues(t,
		1,
		registry1.HistoryTotalRecords(),
	)

	site2 := "Cloudflare"

	req2 := shared.ParamsAddEvent{
		SiteKey: site2,

		IP:              "1.1.1.1",
		Browser:         analytics.Brave,
		OperatingSystem: analytics.Android,

		TimestampUNIX:  localTime.Unix(),
		OffsetUTCHours: int64(utcOffsetHours),

		IsPrivateIP: true,
	}

	req3 := shared.ParamsAddEvent{
		SiteKey: site2,

		IP:              "1.1.1.1",
		Browser:         analytics.Chrome,
		OperatingSystem: analytics.Android,

		TimestampUNIX:  localTime.Unix(),
		OffsetUTCHours: int64(utcOffsetHours),

		IsPrivateIP: true,
	}

	require.Empty(t, dc.AddEvents(&req2, &req3))

	registry2, exists2 := dc.data[Site(site2)]
	require.True(t, exists2)
	require.NotNil(t, registry2)

	require.EqualValues(t,
		0,
		registry2.PreviousMonthTotalRecords(),
	)

	require.EqualValues(t,
		2,
		registry2.CurrentMonthTotalRecords(),
	)
}
