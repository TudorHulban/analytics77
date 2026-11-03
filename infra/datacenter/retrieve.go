package datacenter

import (
	"fmt"
	"io"

	"github.com/tudorhulban/analytics77/domain/analytics"
	"github.com/tudorhulban/analytics77/domain/site"
)

// GetRegistry returns the Registry for the given site.
// If the site does not exist, it returns nil.
//
// This method uses a read lock because it does not mutate the DataCenter.
// The caller must treat the returned pointer as read‑only unless they
// explicitly coordinate ingestion through AddEvents.
func (dc *DataCenter) GetRegistry(siteName site.SiteName) (*analytics.Registry, bool) {
	dc.mu.RLock()

	dcSite, exists := dc.data[siteName]

	dc.mu.RUnlock()

	return dcSite.Data, exists
}

// ListSites returns a snapshot of all sites currently present in the DataCenter.
// The returned slice is a copy of the map keys and is safe for the caller to use.
// This method acquires a read lock and does not mutate the DataCenter.
func (dc *DataCenter) GetSiteNames() []site.SiteName {
	dc.mu.RLock()

	result := make([]site.SiteName, 0, len(dc.data))

	for siteName := range dc.data {
		result = append(result, siteName)
	}

	dc.mu.RUnlock()

	return result
}

// Snapshot writes snapshots of all registries into w.
// It acquires a read lock and delegates snapshotting to each Registry.
func (dc *DataCenter) Snapshot(w io.Writer) error {
	dc.mu.RLock()
	defer dc.mu.RUnlock()

	// Stable ordering of sites
	sites := make([]site.Site, 0, len(dc.data))

	for _, dcSite := range dc.data {
		sites = append(sites, dcSite.Site)
	}

	// slices.Sort(keys)

	for _, item := range sites {
		if _, errPrintSite := fmt.Fprintf(w, "[%s]\n", item); errPrintSite != nil {
			return errPrintSite
		}

		if _, errPrintSnapshot := fmt.Fprintf(w, "[%s]\n", dc.data[item.Name].Site); errPrintSnapshot != nil {
			return errPrintSnapshot
		}

		if _, errPrintEmptyLine := fmt.Fprintln(w); errPrintEmptyLine != nil {
			return errPrintEmptyLine
		}
	}

	return nil
}
