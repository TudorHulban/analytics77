package datacenter

import (
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/tudorhulban/analytics77/domain/analytics"
	"github.com/tudorhulban/hxerrors"
)

// Site is localhost if there is collection from several sites.
type Site string

// Stores in UTC times.
// Converts only at the boundaries.
type DataCenter struct {
	data map[Site]*analytics.Registry
	mu   sync.RWMutex
}

func NewDataCenter(sites ...string) *DataCenter {
	result := DataCenter{
		data: map[Site]*analytics.Registry{},
	}

	for _, site := range sites {
		result.AddSite(site)
	}

	return &result
}

func (dc *DataCenter) AddSite(siteName string) error {
	if len(siteName) == 0 {
		return hxerrors.ErrNilInput{}
	}

	if _, exists := dc.data[Site(siteName)]; !exists {
		dc.data[Site(siteName)] = &analytics.Registry{}
	}

	return nil
}

// Advance is not safe for concurrent callers.
func (dc *DataCenter) Advance(sites ...Site) {
	for site, registry := range dc.data {
		if slices.Contains(sites, site) {
			registry.Advance()
		}
	}
}

func (dc *DataCenter) String() string {
	dc.mu.RLock()
	defer dc.mu.RUnlock()

	var b strings.Builder

	fmt.Fprintf(
		&b,
		"DataCenter: %d registr%s\n",

		len(dc.data),
		func() string {
			if len(dc.data) == 1 {
				return "y"
			}

			return "ies"
		}(),
	)

	keys := make([]Site, 0, len(dc.data))

	for key := range dc.data {
		keys = append(keys, key)
	}

	slices.Sort(keys)

	for _, k := range keys {
		fmt.Fprintf(&b, "\n[%s]\n", k)

		dc.data[k].WriteTo(&b)
	}

	return b.String()
}
