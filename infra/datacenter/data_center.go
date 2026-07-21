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
type DataCenter[T analytics.Ingestable[analytics.TMetric], D analytics.TMetric] struct {
	data map[Site]*analytics.Registry[T, D]
	mu   sync.RWMutex
}

func NewDataCenter[T analytics.Ingestable[analytics.TMetric], D analytics.TMetric]() *DataCenter[T, D] {
	return &DataCenter[T, D]{
		data: map[Site]*analytics.Registry[T, D]{},
	}
}

func (dc *DataCenter[T, D]) String() string {
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

		registryString(dc.data[k], &b)
	}

	return b.String()
}
