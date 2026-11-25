package integrationtests

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/analytics77/domain/analytics"
	"github.com/tudorhulban/analytics77/domain/calendar"
	"github.com/tudorhulban/analytics77/domain/site"
	"github.com/tudorhulban/analytics77/fixtures"
	"github.com/tudorhulban/analytics77/infra/datacenter"
	"github.com/tudorhulban/analytics77/shared"
)

func TestPreviousHourAcrossAprilToMay(t *testing.T) {
	siteName := t.Name()

	// Simulate 2026 May 1, 00:30 UTC.
	may1 := time.Date(
		2026,
		time.May,
		1,
		0,
		30, 0, 0,
		time.UTC,
	)

	dataCenter, errCr := datacenter.NewDataCenter(
		datacenter.ParamsNewDataCenter{
			DefaultIANALocation: fixtures.DefaultLocationIANA,
			DefaultTimeNow: func() time.Time {
				return may1
			},
		},
	)
	require.NoError(t, errCr)
	require.NotNil(t, dataCenter)

	// April 30, 23:30 UTC.
	april30 := may1.Add(-1 * time.Hour)

	// seed datacenter

	// Insert May 1 event.
	reqMay1 := shared.ParamsAddEvent{
		SiteKey: siteName,
		IP:      "127.0.0.1",

		Browser:         analytics.Brave,
		OperatingSystem: analytics.Linux,
		IsPrivateIP:     true,

		TimestampUTC: may1.Unix(),
	}
	require.Empty(t,
		dataCenter.AddEvents(&reqMay1),
	)

	// Insert April 30 event.
	reqApril30 := shared.ParamsAddEvent{
		SiteKey:         siteName,
		IP:              "127.0.0.1",
		Browser:         analytics.Brave,
		OperatingSystem: analytics.Linux,
		IsPrivateIP:     true,

		TimestampUTC: april30.Unix(),
	}
	require.Empty(t,
		dataCenter.AddEvents(&reqApril30),
	)

	// verify seed results

	dataCenter.Snapshot(os.Stdout)

	registry, exists := dataCenter.GetRegistry(site.SiteName(siteName))
	require.True(t, exists)
	require.NotNil(t, registry)

	require.False(t,
		registry.GetActiveSlot().IsZero(),
		"active slot should be valid",
	)
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
		calendar.CalendarDayToIndex(int8(april30.Day())),
	)
	require.NoError(t, errHistory)

	// April 30 slot must contain the event.
	require.False(t,
		aggApril30.IPs.IsZero(),

		"aggregated value is:%s",
		aggApril30.IPs.String(),
	)
}

func TestRolloverAcrossDecemberToJanuary(t *testing.T) {
	siteName := t.Name()

	// Simulate 2026 December 31, 23:30 UTC as current clock.
	dec31 := time.Date(
		2026,
		time.December,
		31,
		23,
		30, 0, 0,
		time.UTC,
	)

	// January 1, 00:30 UTC (next calendar month).
	jan1 := time.Date(
		2027,
		time.January,
		1,
		0,
		30, 0, 0,
		time.UTC,
	)

	dc, errCr := datacenter.NewDataCenter(
		datacenter.ParamsNewDataCenter{
			DefaultIANALocation: fixtures.DefaultLocationIANA,
			DefaultTimeNow: func() time.Time {
				return dec31
			},
		},
	)
	require.NoError(t, errCr)
	require.NotNil(t, dc)

	// Seed December event so registry month is 12.
	errBatch := dc.AddEvents(&shared.ParamsAddEvent{
		SiteKey:         siteName,
		IP:              "127.0.0.1",
		Browser:         analytics.Brave,
		OperatingSystem: analytics.Linux,
		IsPrivateIP:     true,
		TimestampUTC:    dec31.Unix(),
	})
	require.Empty(t, errBatch)

	// Send January event while registry is in December (12).
	// Expected: rollover to January (1) and record the Jan 1 event.
	errBatch = dc.AddEvents(&shared.ParamsAddEvent{
		SiteKey:         siteName,
		IP:              "127.0.0.2",
		Browser:         analytics.Brave,
		OperatingSystem: analytics.Linux,
		IsPrivateIP:     true,
		TimestampUTC:    jan1.Unix(),
	})
	require.Empty(t, errBatch)

	registry, exists := dc.GetRegistry(site.SiteName(siteName))
	require.True(t, exists)
	require.NotNil(t, registry)

	assert.Equal(
		t,
		int32(time.January),
		registry.CalendarMonthCurrentNumber.Load(),
		"registry should roll over from December (12) to January (1)",
	)
	assert.False(
		t,
		registry.GetActiveSlot().IsZero(),
		"active slot should contain the January event after rollover",
	)
}
