package datacenter

import (
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/tudorhulban/analytics77/domain/analytics"
	"github.com/tudorhulban/analytics77/domain/offsets"
	"github.com/tudorhulban/analytics77/domain/site"
	"github.com/tudorhulban/hxerrors"
)

type DCSite struct {
	Data *analytics.Registry
	Site site.Site
}

// Stores in UTC times.
// Converts only at the boundaries.
type DataCenter struct {
	data map[site.SiteName]*DCSite
	mu   sync.RWMutex

	// used in case of automatic site creation.
	defaultOffsets offsets.TimestampOffsets
}

func NewDataCenter(defaultOffsets offsets.TimestampOffsets, sites ...site.Site) (*DataCenter, error) {
	result := DataCenter{
		data: map[site.SiteName]*DCSite{},
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

	if _, exists := dc.data[item.Name]; !exists {
		dc.data[item.Name] = &DCSite{
			Site: item,
			Data: analytics.NewRegistry(
				int8(time.Now().Month()), //nolint:gosec

				item.TimestampDSTSpring,
				item.TimestampDSTWinter,
			),
		}
	}

	return nil
}

// Advance is not safe for concurrent callers.
func (dc *DataCenter) Advance(sites ...site.Site) {
	for _, item := range sites {
		if dcSite, exists := dc.data[item.Name]; exists {
			dcSite.Data.Advance()
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

	for dcSiteName := range dc.data {
		keys = append(keys, dcSiteName)
	}

	slices.Sort(keys)

	for _, siteName := range keys {
		fmt.Fprintf(&b, "\n[%s]\n", siteName)

		fmt.Fprintf(
			&b,
			"\n[%s]\n",
			dc.data[siteName].Site.String(),
		)
	}

	return b.String()
}
