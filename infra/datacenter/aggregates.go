package datacenter

import (
	"fmt"

	"github.com/tudorhulban/analytics77/domain/analytics"
)

func (dc *DataCenter) CurrentMonthAggregateTopNForHour(day, hour int8, sites ...Site) ([]*analytics.AggregatedTopN, error) {
	if len(sites) == 0 {
		return nil, analytics.ErrInvalidInput
	}

	result := make([]*analytics.AggregatedTopN, 0)

	for _, site := range sites {
		if registry, exists := dc.GetRegistry(site); exists {
			aggregate, errAgg := registry.CurrentMonthAggregateTopNForHour(day, hour)
			if errAgg != nil {
				return nil,
					fmt.Errorf(
						"error for site: %q : %w",
						site,
						errAgg,
					)
			}

			result = append(result, aggregate)
		}
	}

	return result, nil
}

func (dc *DataCenter) CurrentMonthAggregateTopN(sites ...Site) ([]*analytics.AggregatedTopN, error) {
	if len(sites) == 0 {
		return nil, analytics.ErrInvalidInput
	}

	result := make([]*analytics.AggregatedTopN, 0)

	for _, site := range sites {
		if registry, exists := dc.GetRegistry(site); exists {
			aggregate := registry.CurrentMonthAggregateTopN()

			if aggregate != nil {
				result = append(result, aggregate)
			}
		}
	}

	return result, nil
}
