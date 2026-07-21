package datacenter

import (
	"fmt"
	"io"
	"slices"

	"github.com/tudorhulban/analytics77/domain/analytics"
)

// GetRegistry returns the Registry for the given site.
// If the site does not exist, it returns nil.
//
// This method uses a read lock because it does not mutate the DataCenter.
// The caller must treat the returned pointer as read‑only unless they
// explicitly coordinate ingestion through AddEvents.
func (dc *DataCenter[T, D]) GetRegistry(site Site) *analytics.Registry[T, D] {
	dc.mu.RLock()

	result := dc.data[site]

	dc.mu.RUnlock()

	return result
}

// ListSites returns a snapshot of all sites currently present in the DataCenter.
// The returned slice is a copy of the map keys and is safe for the caller to use.
// This method acquires a read lock and does not mutate the DataCenter.
func (dc *DataCenter[T, D]) GetSiteNames() []Site {
	dc.mu.RLock()

	result := make([]Site, 0, len(dc.data))
	for site := range dc.data {
		result = append(result, site)
	}

	dc.mu.RUnlock()

	return result
}

// Snapshot writes snapshots of all registries into w.
// It acquires a read lock and delegates snapshotting to each Registry.
func (dc *DataCenter[T, D]) Snapshot(w io.Writer) error {
	dc.mu.RLock()
	defer dc.mu.RUnlock()

	// Stable ordering of sites
	keys := make([]string, 0, len(dc.data))

	for dcSite := range dc.data {
		keys = append(keys, string(dcSite))
	}

	slices.Sort(keys)

	for _, site := range keys {
		if _, errPrintSite := fmt.Fprintf(w, "[%s]\n", site); errPrintSite != nil {
			return errPrintSite
		}

		if errPrintSnapshot := dc.data[Site(site)].Snapshot(w); errPrintSnapshot != nil {
			return errPrintSnapshot
		}

		if _, errPrintEmptyLine := fmt.Fprintln(w); errPrintEmptyLine != nil {
			return errPrintEmptyLine
		}
	}

	return nil
}
