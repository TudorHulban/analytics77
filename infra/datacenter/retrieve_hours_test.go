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

func TestCurrentDayHoursWithData(t *testing.T) {
	offsetsROU := helpers.TimestampOffsets{
		OffsetUTCHours: +3,
	}

	dc := NewDataCenter()
	require.NotNil(t, dc)

	hoursWithDataBefore, howManyBefore := dc.CurrentDayHoursWithData(&offsetsROU)
	require.Zero(t, howManyBefore)
	require.Zero(t, hoursWithDataBefore)

	site1 := "localhost"
	localTime := time.Now()

	recordsPerSiteBefore := dc.GetPreviousHourRecordsPerSite(&offsetsROU)
	require.Empty(t, recordsPerSiteBefore)

	request1 := shared.ParamsAddEvent{
		SiteKey: site1,

		IP:              "127.0.0.1",
		Browser:         analytics.Brave,
		OperatingSystem: analytics.Linux,

		TimestampUNIX: helpers.GetHourInOffsetAsUTC(
			localTime,
			localTime.Hour(),
			int(offsetsROU.OffsetUTCHours),
		),

		IsPrivateIP: true,
	}

	require.Empty(t, dc.AddEvents(&request1))

	hoursWithDataAfter1, howManyAfter1 := dc.CurrentDayHoursWithData(&offsetsROU)
	require.EqualValues(t, 1, howManyAfter1, "howManyAfter should be 1")
	require.NotZero(t, hoursWithDataAfter1, "hoursWithDataAfter should be not zero in most cases")
	require.EqualValues(t,
		localTime.Hour()-int(offsetsROU.OffsetUTCHours),
		hoursWithDataAfter1[0],
	)

	recordsPerSiteAfter1 := dc.GetPreviousHourRecordsPerSite(&offsetsROU)
	require.Len(t, recordsPerSiteAfter1, 1, "should be 1 as site was added")
	require.Zero(t, recordsPerSiteAfter1[Site(site1)])

	oneHourBefore := localTime.Add(-1 * time.Hour)

	request2 := shared.ParamsAddEvent{
		SiteKey: site1,

		IP:              "127.0.0.1",
		Browser:         analytics.Brave,
		OperatingSystem: analytics.Linux,

		TimestampUNIX: helpers.GetHourInOffsetAsUTC(oneHourBefore, oneHourBefore.Hour(), int(offsetsROU.OffsetUTCHours)),

		IsPrivateIP: true,
	}

	require.Empty(t, dc.AddEvents(&request2))

	hoursWithDataAfter2, howManyAfter2 := dc.CurrentDayHoursWithData(&offsetsROU)
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
