package integrationtests

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/analytics77/domain/analytics"
	"github.com/tudorhulban/analytics77/domain/dhelpers"
	"github.com/tudorhulban/analytics77/infra/datacenter"
	"github.com/tudorhulban/analytics77/shared"
)

func TestPreviousHourAcrossAprilToMay(t *testing.T) {
	dc := datacenter.NewDataCenter()
	require.NotNil(t, dc)

	site := "localhost"

	// Simulate 2026 May 1, 00:30 UTC.
	may1 := time.Date(2026, time.May, 1, 0, 30, 0, 0, time.UTC)

	// April 30, 23:30 UTC.
	april30 := may1.Add(-1 * time.Hour)

	// Insert May 1 event.
	reqMay1 := shared.ParamsAddEvent{
		SiteKey: site,
		IP:      "127.0.0.1",

		Browser:         analytics.Brave,
		OperatingSystem: analytics.Linux,

		TimestampUNIX: may1.Unix(),

		IsPrivateIP: true,
	}
	require.Empty(t, dc.AddEvents(&reqMay1))

	// Insert April 30 event.
	reqApril30 := shared.ParamsAddEvent{
		SiteKey:         site,
		IP:              "127.0.0.1",
		Browser:         analytics.Brave,
		OperatingSystem: analytics.Linux,

		TimestampUNIX: april30.Unix(),

		IsPrivateIP: true,
	}
	require.Empty(t, dc.AddEvents(&reqApril30))

	registry := dc.GetRegistry(datacenter.Site(site))
	require.False(t, registry.GetActiveSlot().IsZero())
	require.False(t, registry.GetPreviousSlot().IsZero())

	var bufBefore strings.Builder

	registry.Snapshot(&bufBefore)
	t.Log("registry before advance:\n", bufBefore.String())

	registry.Advance() // move April in history as month 0.

	var bufAfter strings.Builder

	registry.Snapshot(&bufAfter)
	t.Log("registry after advance:\n", bufAfter.String())

	// Query April 30 aggregated data.
	aggApril30, errHistory := registry.HistoryAggregateTopNForDay(
		analytics.FromTwoMonthsAgo,
		dhelpers.CalendarDayToIndex(int8(april30.Day())),
	)
	require.NoError(t, errHistory)

	// April 30 slot must contain the event.
	require.False(t,
		aggApril30.IPs.IsZero(),

		"aggregated value is:%s",
		aggApril30.IPs.String(),
	)
}
