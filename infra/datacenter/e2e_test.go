package datacenter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/analytics77/domain/analytics"
	"github.com/tudorhulban/analytics77/domain/site"
	"github.com/tudorhulban/analytics77/shared"
)

func TestDataCenter(t *testing.T) {
	dc, errCr := NewDataCenter()
	require.NoError(t, errCr)
	require.NotNil(t, dc)

	require.Empty(t, dc.AddEvents())
	require.Empty(t, dc.AddEvents(nil))

	site1Test := site.Site{
		Name: "xxxx - " + site.SiteName(t.Name()),
	}

	localTime := time.Now()
	_, offsetSeconds := localTime.Zone()

	utcOffsetHours := offsetSeconds / 3600

	req1 := shared.ParamsAddEvent{
		SiteKey: string(site1Test.Name),

		IP:              "127.0.0.1",
		Browser:         analytics.Brave,
		OperatingSystem: analytics.Linux,

		TimestampUNIX:  localTime.Unix(),
		OffsetUTCHours: int64(utcOffsetHours),

		IsPrivateIP: true,
	}

	require.Empty(t, dc.AddEvents(&req1))

	require.Len(t, dc.data, 1, "site1 should have been added")

	dsSite1, exists1 := dc.data[site1Test.Name]
	require.True(t, exists1)
	require.NotNil(t, dsSite1)

	siteNames := dc.GetSiteNames()
	require.Len(t, siteNames, 1)
	require.EqualValues(t, req1.SiteKey, siteNames[0])

	require.False(t,
		dsSite1.
			Data.
			GetActiveSlot().
			IsZero(),
	)
	require.True(t,
		dsSite1.
			Data.
			GetPreviousSlot().
			IsZero(),
	)
	require.True(t,
		dsSite1.
			Data.
			HistoryAggregateTopN().
			IsZero(),
	)

	hoursWithData, howMany := dsSite1.
		Data.
		CurrentMonthHoursWithData(int8(localTime.Day()))
	require.EqualValues(t, 1, howMany)
	require.EqualValues(t,
		localTime.Hour(),
		hoursWithData[0],
	)

	require.EqualValues(t,
		0,
		dsSite1.
			Data.
			PreviousMonthTotalRecords(),
	)

	require.EqualValues(t,
		1,
		dsSite1.
			Data.
			CurrentMonthTotalRecords(),
	)

	aggregates1 := dsSite1.Data.CurrentMonthAggregateTopN()
	require.EqualValues(t,
		analytics.Brave,
		*aggregates1.Browsers.Names[0].Load(),

		"browser want:%s got:%s",
		analytics.Brave,
		aggregates1.Browsers.Names[0].Load(),
	)

	dc.Advance(site1Test)

	require.EqualValues(t,
		0,
		dsSite1.
			Data.
			CurrentMonthTotalRecords(),
	)

	require.EqualValues(t,
		1,
		dsSite1.
			Data.
			PreviousMonthTotalRecords(),
	)

	dc.Advance(site1Test)

	require.EqualValues(t,
		1,
		dsSite1.
			Data.
			HistoryTotalRecords(),
	)

	site2Test := site.Site{
		Name: "yyyy - " + site.SiteName(t.Name()),
	}

	req2 := shared.ParamsAddEvent{
		SiteKey: string(site2Test.Name),

		IP:              "1.1.1.1",
		Browser:         analytics.Brave,
		OperatingSystem: analytics.Android,

		TimestampUNIX:  localTime.Unix(),
		OffsetUTCHours: int64(utcOffsetHours),

		IsPrivateIP: true,
	}

	req3 := shared.ParamsAddEvent{
		SiteKey: string(site2Test.Name),

		IP:              "1.1.1.1",
		Browser:         analytics.Chrome,
		OperatingSystem: analytics.Android,

		TimestampUNIX:  localTime.Unix(),
		OffsetUTCHours: int64(utcOffsetHours),

		IsPrivateIP: true,
	}

	require.Empty(t, dc.AddEvents(&req2, &req3))

	dcSite2, exists2 := dc.data[site2Test.Name]
	require.True(t, exists2)
	require.NotNil(t, dcSite2)

	require.EqualValues(t,
		0,
		dcSite2.
			Data.
			PreviousMonthTotalRecords(),
	)

	require.EqualValues(t,
		2,
		dcSite2.
			Data.
			CurrentMonthTotalRecords(),
	)
}
