package datacenter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/analytics77/domain/analytics"
	"github.com/tudorhulban/analytics77/domain/offsets"
	"github.com/tudorhulban/analytics77/domain/site"
	"github.com/tudorhulban/analytics77/shared"
)

func TestDataCenter(t *testing.T) {
	datacenter, errCrDatacenter := NewDataCenter()
	require.NoError(t, errCrDatacenter)
	require.NotNil(t, datacenter)

	require.Empty(t, datacenter.AddEvents())
	require.Empty(t, datacenter.AddEvents(nil))

	site1Test := site.Site{
		Name: "xxxx - " + site.SiteName(t.Name()),

		TimestampOffsets: offsets.TimestampOffsets{
			OffsetUTCHours: -2,
		},
	}

	localTime := time.Now()

	req1 := shared.ParamsAddEvent{
		SiteKey: string(site1Test.Name),

		IP:              "127.0.0.1",
		Browser:         analytics.Brave,
		OperatingSystem: analytics.Linux,

		TimestampUNIX: localTime.Unix(),
		IsPrivateIP:   true,
	}

	require.Empty(t, datacenter.AddEvents(&req1))

	require.Len(t, datacenter.data, 1, "site1 should have been added")

	dcSite1, exists1 := datacenter.data[site1Test.Name]
	require.True(t, exists1)
	require.NotNil(t, dcSite1)
	require.Equal(t,
		site1Test.Name,
		dcSite1.Site.Name,
	)

	siteNames := datacenter.GetSiteNames()
	require.Len(t, siteNames, 1)
	require.EqualValues(t, req1.SiteKey, siteNames[0])

	require.False(t,
		dcSite1.
			Data.
			GetActiveSlot().
			IsZero(),
	)
	require.True(t,
		dcSite1.
			Data.
			GetPreviousSlot().
			IsZero(),
	)
	require.True(t,
		dcSite1.
			Data.
			HistoryAggregateTopN().
			IsZero(),
	)

	// dcSite1.Data.Snapshot(os.Stdout)

	hoursWithData, howMany := dcSite1.
		Data.
		CurrentMonthHoursWithData(int8(localTime.Day()))
	require.EqualValues(t, 1, howMany)
	require.EqualValues(t,
		localTime.Hour()+int(site1Test.OffsetUTCHours),
		hoursWithData[0],
	)

	require.EqualValues(t,
		0,
		dcSite1.
			Data.
			PreviousMonthTotalRecords(),
	)

	require.EqualValues(t,
		1,
		dcSite1.
			Data.
			CurrentMonthTotalRecords(),
	)

	aggregates1 := dcSite1.Data.CurrentMonthAggregateTopN()
	require.EqualValues(t,
		analytics.Brave,
		*aggregates1.Browsers.Names[0].Load(),

		"browser want:%s got:%s",
		analytics.Brave,
		aggregates1.Browsers.Names[0].Load(),
	)

	datacenter.Advance(site1Test)

	require.EqualValues(t,
		0,
		dcSite1.
			Data.
			CurrentMonthTotalRecords(),
	)

	require.EqualValues(t,
		1,
		dcSite1.
			Data.
			PreviousMonthTotalRecords(),
	)

	datacenter.Advance(site1Test)

	require.EqualValues(t,
		1,
		dcSite1.
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

		TimestampUNIX: localTime.Unix(),
		IsPrivateIP:   true,
	}

	req3 := shared.ParamsAddEvent{
		SiteKey: string(site2Test.Name),

		IP:              "1.1.1.1",
		Browser:         analytics.Chrome,
		OperatingSystem: analytics.Android,

		TimestampUNIX: localTime.Unix(),
		IsPrivateIP:   true,
	}

	require.Empty(t, datacenter.AddEvents(&req2, &req3))

	dcSite2, exists2 := datacenter.data[site2Test.Name]
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
