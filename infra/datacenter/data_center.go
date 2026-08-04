package datacenter

import (
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/tudorhulban/analytics77/domain/analytics"
)

type Site string

// Stores in UTC times.
// Converts only at the boundaries.
type DataCenter struct {
	data map[Site]*analytics.Registry
	mu   sync.RWMutex
}

func NewDataCenter() *DataCenter {
	return &DataCenter{
		data: map[Site]*analytics.Registry{},
	}
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
