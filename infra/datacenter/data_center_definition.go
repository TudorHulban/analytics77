package datacenter

import (
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/tudorhulban/analytics77/domain/analytics"
	"github.com/tudorhulban/analytics77/domain/site"
	"github.com/tudorhulban/hxerrors"
)

// Stores in UTC times.
// Converts only at the boundaries.
type DataCenter struct {
	data map[site.Site]*analytics.Registry
	mu   sync.RWMutex
}

func NewDataCenter(sites ...site.Site) (*DataCenter, error) {
	result := DataCenter{
		data: map[site.Site]*analytics.Registry{},
	}

	for _, item := range sites {
		if errAdd := result.AddSite(item); errAdd != nil {
			return nil,
				errAdd
		}
	}

	return &result,
		nil
}

func (dc *DataCenter) AddSite(item site.Site) error {
	if len(item.Name) == 0 {
		return hxerrors.ErrInvalidInput{
			Issue:      hxerrors.ErrNilInput{},
			InputValue: item,
			InputName:  "site",
			Caller:     "AddSite",
		}
	}

	if _, exists := dc.data[item]; !exists {
		dc.data[item] = &analytics.Registry{}
	}

	return nil
}

// Advance is not safe for concurrent callers.
func (dc *DataCenter) Advance(sites ...site.Site) {
	for siteName, registry := range dc.data {
		if slices.Contains(sites, siteName) {
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

	keys := make([]site.SiteName, 0, len(dc.data))

	for item := range dc.data {
		keys = append(keys, item.Name)
	}

	slices.Sort(keys)

	for _, siteName := range keys {
		fmt.Fprintf(&b, "\n[%s]\n", siteName)

		// TODO: needs review
		dc.data[siteName].WriteTo(&b)
	}

	return b.String()
}
