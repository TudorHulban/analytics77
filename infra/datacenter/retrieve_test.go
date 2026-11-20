package datacenter

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/analytics77/domain/site"
	"github.com/tudorhulban/analytics77/fixtures"

	"github.com/tudorhulban/analytics77/shared"
)

func TestSnapshot(t *testing.T) {
	numSites := 3
	// recordsPerSite := 5

	sites := make([]site.Site, numSites)

	for ix := range numSites {
		sites[ix] = site.Site{
			Name:         site.SiteName(fmt.Sprintf("site-%d", ix)),
			SiteLocation: fixtures.DefaultLocationIANA,
		}
	}

	datacenter, errCrDatacenter := NewDataCenter(
		fixtures.DefaultLocationIANA,
		sites...,
	)
	require.NoError(t, errCrDatacenter)

	// Execute Snapshot
	var buf bytes.Buffer

	require.NoError(t,
		datacenter.Snapshot(&buf),
	)

	output := buf.String()

	fmt.Println(output)

	for ix := range numSites {
		require.Contains(t,
			output,
			fmt.Sprintf("site-%d", ix),
		)
	}
}

func TestGetPreviousHourRecordsPerSite_CrossesMonthBoundary(t *testing.T) {
	datacenter, errCr := NewDataCenter(fixtures.DefaultLocationIANA)
	require.NoError(t, errCr)
	require.NotNil(t, datacenter)

	siteName := "localhost"

	// 2024 is a leap year: Feb has 29 days. If the code just swapped a
	// hardcoded "31" for a hardcoded "28/29/30", this would still fail.
	lastHourFebUTC := time.Date(
		2024,
		time.February,
		29,
		23, 15, 0, 0,
		time.UTC,
	)
	nowMarchUTC := time.Date(
		2024,
		time.March,
		1,
		0, 30, 0, 0,
		time.UTC,
	)

	// First event seeds the registry's current month as February.
	require.Empty(t,
		datacenter.AddEvents(
			&shared.ParamsAddEvent{
				SiteKey: siteName,
				IP:      "127.0.0.1",

				TimestampUTC: lastHourFebUTC.Unix(),
				IsPrivateIP:  true,
			},
		),
	)

	// Second event is one calendar month later -> triggers AddEvents' own
	// Rollover() path, pushing February into "previous" and March into
	// "current". This mirrors how the registry actually gets here in
	// production; it is not test-only setup.
	require.Empty(t,
		datacenter.AddEvents(
			&shared.ParamsAddEvent{
				SiteKey: siteName,
				IP:      "127.0.0.1",

				TimestampUTC: nowMarchUTC.Unix(),
				IsPrivateIP:  true,
			},
		),
	)

	registry, exists := datacenter.GetRegistry(site.SiteName(siteName))
	require.True(t, exists)

	require.NotNil(t, registry)
	require.EqualValues(t,
		time.March,
		registry.CalendarMonthCurrentNumber.Load(),
	)

	// At March 1st 00:30, "the previous hour" is Feb 29 23:xx -> the data
	// that belongs to it lives in GetPreviousMonth(), not GetCurrentMonth().
	records := datacenter.PreviousHourRecordsPerSiteAt(nowMarchUTC)

	require.Len(t, records, 1)
	require.EqualValues(t,
		1,
		records[site.SiteName(siteName)],

		"the last hour of February should be visible as 'previous hour' from March 1st 00:30",
	)
}
