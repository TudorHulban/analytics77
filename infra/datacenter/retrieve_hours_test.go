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

func TestCurrentDayHoursWithData(t *testing.T) {
	offsetsROU := offsets.TimestampOffsets{
		OffsetUTCHours: +3,
	}

	datacenter, errCr := NewDataCenter(
		ParamsNewDataCenter{
			DefaultIANALocation: fixtures.DefaultLocationIANA,
		},
	)
	require.NoError(t, errCr)
	require.NotNil(t, datacenter)

	hoursWithDataBefore, howManyBefore := datacenter.CurrentDayHoursWithData(&offsetsROU)
	require.Zero(t, howManyBefore)
	require.Zero(t, hoursWithDataBefore)

	site1Test := site.Site{
		Name: "xxxx - " + site.SiteName(t.Name()),
	}
	localTime := time.Now()

	recordsPerSiteBefore := datacenter.GetPreviousHourRecordsPerSite()
	require.Empty(t, recordsPerSiteBefore)

	request1 := shared.ParamsAddEvent{
		SiteKey: string(site1Test.Name),

		IP:              "127.0.0.1",
		Browser:         analytics.Brave,
		OperatingSystem: analytics.Linux,

		TimestampUTC: offsets.GetHourInOffsetAsUTC(
			localTime,
			localTime.Hour(),
			int(offsetsROU.OffsetUTCHours),
		),

		IsPrivateIP: true,
	}

	require.Empty(t, datacenter.AddEvents(&request1))

	hoursWithDataAfter1, howManyAfter1 := datacenter.CurrentDayHoursWithData(&offsetsROU)
	require.EqualValues(t,
		1,
		howManyAfter1,
		"howManyAfter should be 1",
	)
	require.NotZero(t,
		hoursWithDataAfter1,
		"hoursWithDataAfter should be not zero in most cases",
	)
	require.EqualValues(t,
		localTime.Hour()-int(offsetsROU.OffsetUTCHours),
		hoursWithDataAfter1[0],
	)

	recordsPerSiteAfter1 := datacenter.GetPreviousHourRecordsPerSite()
	require.Len(t,
		recordsPerSiteAfter1,
		1,
		"should be 1 as site was added",
	)
	require.Zero(t,
		recordsPerSiteAfter1[site1Test.Name],
	)

	oneHourBefore := localTime.Add(-1 * time.Hour)

	request2 := shared.ParamsAddEvent{
		SiteKey: string(site1Test.Name),

		IP:              "127.0.0.1",
		Browser:         analytics.Brave,
		OperatingSystem: analytics.Linux,

		TimestampUTC: offsets.GetHourInOffsetAsUTC(oneHourBefore, oneHourBefore.Hour(), int(offsetsROU.OffsetUTCHours)),

		IsPrivateIP: true,
	}

	require.Empty(t, datacenter.AddEvents(&request2))

	hoursWithDataAfter2, howManyAfter2 := datacenter.CurrentDayHoursWithData(&offsetsROU)
	require.EqualValues(t,
		2,
		howManyAfter2,
		"howManyAfter should be 1",
	)
	require.NotZero(t,
		hoursWithDataAfter2,
		"hoursWithDataAfter should be not zero in most cases",
	)
	require.Empty(t,
		hxhelpers.NotInSliceSource(hoursWithDataAfter2[:howManyAfter2], int8(oneHourBefore.Hour()-int(offsetsROU.OffsetUTCHours)), int8(localTime.Hour()-int(offsetsROU.OffsetUTCHours))),
		hoursWithDataAfter2,
	)
}
