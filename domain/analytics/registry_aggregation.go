package analytics

func (*Registry) aggregateHour(fromMetric *MetricActive) *AggregatedTopN {
	var result AggregatedTopN

	if (*fromMetric).GetRecordsPerPeriod() == 0 {
		return &result
	}

	(*fromMetric).GetTopIPs().DeepCopyInto(&result.IPs)
	(*fromMetric).GetTopASNs().DeepCopyInto(&result.ASN)
	(*fromMetric).GetTopCountries().DeepCopyInto(&result.Countries)
	(*fromMetric).GetTopCities().DeepCopyInto(&result.Cities)
	(*fromMetric).GetTopURLs().DeepCopyInto(&result.URL)

	(*fromMetric).GetTopOperatingSystems().DeepCopyInto(&result.OS)
	(*fromMetric).GetTopBrowsers().DeepCopyInto(&result.Browsers)

	return &result
}

func (r *Registry) CurrentMonthAggregateTopNForHour(day, hour int8) (*AggregatedTopN, error) {
	if day < 0 || day >= 31 || hour < 0 || hour >= 24 {
		return nil,
			ErrInvalidInput
	}

	currentMonth := r.GetActiveSlot()

	return r.aggregateHour((*currentMonth).GetMetric(day, hour)), nil
}

func (r *Registry) PreviousMonthAggregateTopNForHour(day, hour int8) (*AggregatedTopN, error) {
	if day < 0 || day >= 31 || hour < 0 || hour >= 24 {
		return nil,
			ErrInvalidInput
	}

	previousMonth := r.GetPreviousSlot()

	return r.aggregateHour((*previousMonth).GetMetric(day, hour)), nil
}

func (r *Registry) PreviousMonthTotalRecords() uint32 {
	var result uint32

	previousMonth := r.GetPreviousSlot()

	for day := range int8(31) {
		for hour := range int8(24) {
			result = result + (*previousMonth).GetMetric(day, hour).GetRecordsPerPeriod()
		}
	}

	return result
}

func (r *Registry) PreviousMonthTotalRecordsForDay(day int8) uint32 {
	if day < 0 || day >= 31 {
		return 0
	}

	previousSlot := r.GetPreviousSlot()

	var result uint32

	r.ForEachMetric(
		previousSlot,
		func(d int8, h int8, m *MetricActive) {
			if d == day {
				result = result + (*m).GetRecordsPerPeriod()
			}
		},
	)

	return result
}

func (r *Registry) CurrentMonthTotalRecords() uint32 {
	var result uint32

	currentSlot := r.GetActiveSlot()

	for day := range int8(31) {
		for hour := range int8(24) {
			result = result + (*currentSlot).GetMetric(day, hour).GetRecordsPerPeriod()
		}
	}

	return result
}

func (r *Registry) CurrentMonthTotalRecordsForDay(day int8) uint32 {
	var result uint32

	if day < 0 || day >= 31 {
		return result
	}

	currentSlot := r.GetActiveSlot()

	for hour := range int8(24) {
		m := (*currentSlot).GetMetric(day, hour)

		result = result + (*m).GetRecordsPerPeriod()
	}

	return result
}

func (*Registry) mergeHourInto(from *MetricActive, dst *AggregatedTopN) {
	if from.GetRecordsPerPeriod() == 0 {
		return
	}

	dst.IPs.MergeFrom(from.GetTopIPs())
	dst.ASN.MergeFrom(from.GetTopASNs())
	dst.Countries.MergeFrom(from.GetTopCountries())
	dst.Cities.MergeFrom(from.GetTopCities())
	dst.URL.MergeFrom(from.GetTopURLs())

	dst.OS.MergeFrom(from.GetTopOperatingSystems())
	dst.Browsers.MergeFrom(from.GetTopBrowsers())
}

func (r *Registry) PreviousMonthAggregateTopNForDay(day int8) (*AggregatedTopN, error) {
	if day < 0 || day >= 31 {
		return nil,
			ErrInvalidInput
	}

	var result AggregatedTopN

	previousSlot := r.GetPreviousSlot()

	for hour := range int8(24) {
		r.mergeHourInto((*previousSlot).GetMetric(day, hour), &result)
	}

	return &result, nil
}

// PreviousMonthAggregateTopN aggregates Top‑N metrics
// across the entire previous month.
// It traverses every hour of the previous month
// and merges its Top‑N counters into a single cumulative Top‑N structure.
// It accumulates all Top‑N metrics across the entire previous month.
func (r *Registry) PreviousMonthAggregateTopN() *AggregatedTopN {
	var result AggregatedTopN

	r.PreviousMonthForEach(
		func(_, _ int8, m *MetricActive) {
			r.mergeHourInto(m, &result)
		},
	)

	return &result
}

func (r *Registry) CurrentMonthAggregateTopNForDay(day int8) (*AggregatedTopN, error) {
	if day < 0 || day >= 31 {
		return nil,
			ErrInvalidInput
	}

	var result AggregatedTopN

	currentSlot := r.GetActiveSlot()

	for hour := range int8(24) {
		r.mergeHourInto((*currentSlot).GetMetric(day, hour), &result)
	}

	return &result, nil
}

func (r *Registry) CurrentMonthAggregateTopN() *AggregatedTopN {
	var result AggregatedTopN

	r.CurrentMonthForEach(
		func(_, _ int8, m *MetricActive) {
			r.mergeHourInto(m, &result)
		},
	)

	return &result
}

func (r *Registry) HistoryAggregateTopNForDay(month MonthsBack, day int8) (*AggregatedTopN, error) {
	if month >= 6 || day < 0 || day >= 31 {
		return nil, ErrInvalidInput
	}

	var result AggregatedTopN

	slotHistory, errHistory := r.GetHistorySlot(month)
	if errHistory != nil {
		return nil,
			errHistory
	}

	for hour := range int8(24) {
		m := slotHistory.GetMetric(day, hour)

		if m.GetRecordsPerPeriod() == 0 {
			continue
		}

		result.IPs.MergeFrom(m.GetTopIPs())
		result.ASN.MergeFrom(m.GetTopASNs())
		result.Countries.MergeFrom(m.GetTopCountries())
		result.Cities.MergeFrom(m.GetTopCities())
		result.URL.MergeFrom(m.GetTopURLs())
		result.OS.MergeFrom(m.GetTopOperatingSystems())
		result.Browsers.MergeFrom(m.GetTopBrowsers())
	}

	return &result, nil
}

func (r *Registry) HistoryAggregateTopNForMonth(month MonthsBack) (*AggregatedTopN, error) {
	if month >= 6 {
		return nil,
			ErrInvalidInput
	}

	var result AggregatedTopN

	slotHistory, errHistory := r.GetHistorySlot(month)
	if errHistory != nil {
		return nil,
			errHistory
	}

	for day := range int8(31) {
		for hour := range int8(24) {
			fromMetric := (*slotHistory).GetMetric(day, hour)

			records := (*fromMetric).GetRecordsPerPeriod()
			if records == 0 {
				continue
			}

			// Aggregate IPs
			topIPs := (*fromMetric).GetTopIPs()
			for i := range topIPs.Names {
				namePtr := topIPs.Names[i].Load()
				if namePtr == nil {
					continue
				}

				val := topIPs.Values[i].Load()
				if val > 0 {
					result.IPs.Increment(*namePtr, val)
				}
			}

			// Aggregate ASN
			topASN := (*fromMetric).GetTopASNs()
			for i := range topASN.Names {
				namePtr := topASN.Names[i].Load()
				if namePtr == nil {
					continue
				}

				val := topASN.Values[i].Load()
				if val > 0 {
					result.ASN.Increment(*namePtr, val)
				}
			}

			// Aggregate Countries
			topCountries := (*fromMetric).GetTopCountries()
			for i := range topCountries.Names {
				namePtr := topCountries.Names[i].Load()
				if namePtr == nil {
					continue
				}

				val := topCountries.Values[i].Load()
				if val > 0 {
					result.Countries.Increment(*namePtr, val)
				}
			}

			// Aggregate Cities
			topCities := (*fromMetric).GetTopCities()
			for i := range topCities.Names {
				namePtr := topCities.Names[i].Load()
				if namePtr == nil {
					continue
				}

				val := topCities.Values[i].Load()
				if val > 0 {
					result.Cities.Increment(*namePtr, val)
				}
			}

			// Aggregate URLs
			topURLs := (*fromMetric).GetTopURLs()
			for i := range topURLs.Names {
				namePtr := topURLs.Names[i].Load()
				if namePtr == nil {
					continue
				}

				val := topURLs.Values[i].Load()
				if val > 0 {
					result.URL.Increment(*namePtr, val)
				}
			}

			// Aggregate Operating Systems
			topOS := (*fromMetric).GetTopOperatingSystems()
			for i := range topOS.Names {
				namePtr := topOS.Names[i].Load()
				if namePtr == nil {
					continue
				}

				val := topOS.Values[i].Load()
				if val > 0 {
					result.OS.Increment(*namePtr, val)
				}
			}

			// Aggregate Browsers
			topBrowsers := (*fromMetric).GetTopBrowsers()
			for i := range topBrowsers.Names {
				namePtr := topBrowsers.Names[i].Load()
				if namePtr == nil {
					continue
				}

				val := topBrowsers.Values[i].Load()
				if val > 0 {
					result.Browsers.Increment(*namePtr, val)
				}
			}
		}
	}

	return &result, nil
}

func (r *Registry) HistoryAggregateTopNForHour(month MonthsBack, day, hour int8) (*AggregatedTopN, error) {
	if month >= 6 || day < 0 || day >= 31 || hour < 0 || hour >= 24 {
		return nil,
			ErrInvalidInput
	}

	var result AggregatedTopN

	slotHistory, errHistory := r.GetHistorySlot(month)
	if errHistory != nil {
		return nil,
			errHistory
	}

	fromMetric := (*slotHistory).GetMetric(day, hour)

	if (*fromMetric).GetRecordsPerPeriod() == 0 {
		return &result, nil
	}

	(*fromMetric).GetTopIPs().DeepCopyInto(&result.IPs)
	(*fromMetric).GetTopASNs().DeepCopyInto(&result.ASN)
	(*fromMetric).GetTopCountries().DeepCopyInto(&result.Countries)
	(*fromMetric).GetTopCities().DeepCopyInto(&result.Cities)
	(*fromMetric).GetTopURLs().DeepCopyInto(&result.URL)

	(*fromMetric).GetTopOperatingSystems().DeepCopyInto(&result.OS)
	(*fromMetric).GetTopBrowsers().DeepCopyInto(&result.Browsers)

	return &result, nil
}

func (r *Registry) HistoryAggregateTopN() *AggregatedTopN {
	var result AggregatedTopN

	// iterate all history months (0 = previous, 1 = two months ago, ...)
	for months := range uint8(len(r.Slots) - 1) {
		slot, errGetHistorySlot := r.GetHistorySlot(MonthsBack(months))
		if errGetHistorySlot != nil {
			continue
		}

		// reuse the existing aggregator for a single month
		r.MonthForEach(slot,
			func(_, _ int8, m *MetricActive) {
				r.mergeHourInto(m, &result)
			},
		)
	}

	return &result
}

func (r *Registry) HistoryTotalRecords() uint32 {
	var result uint32

	for months := range uint8(len(r.Slots) - 1) {
		slot, errGetHistorySlot := r.GetHistorySlot(MonthsBack(months))
		if errGetHistorySlot != nil {
			continue
		}

		for day := range int8(31) {
			for hour := range int8(24) {
				result = result + (*slot).GetMetric(day, hour).GetRecordsPerPeriod()
			}
		}
	}

	return result
}

// HistoryTotalRecordsForMonth returns the total RecordsPerPeriod for the
// logical history month specified by `monthsBack` (0..5). The function resolves
// the corresponding physical slot via GetHistorySlot and sums only archived
// metrics. The current slot is never included.
func (r *Registry) HistoryTotalRecordsForMonth(monthsBack MonthsBack) uint32 {
	if monthsBack >= 6 {
		return 0
	}

	slot, errGetHistorySlot := r.GetHistorySlot(monthsBack)
	if errGetHistorySlot != nil {
		return 0
	}

	var result uint32

	for day := int8(0); day < 31; day++ {
		for hour := int8(0); hour < 24; hour++ {
			result = result + slot.GetMetric(day, hour).GetRecordsPerPeriod()
		}
	}

	return result
}

func (r *Registry) HistoryTotalRecordsForDay(monthsBack MonthsBack, day int8) uint32 {
	if monthsBack >= 6 || day < 0 || day >= 31 {
		return 0
	}

	slot, errGetHistorySlot := r.GetHistorySlot(monthsBack)
	if errGetHistorySlot != nil {
		return 0
	}

	var result uint32

	for hour := range int8(24) {
		result = result + slot.GetMetric(day, hour).GetRecordsPerPeriod()
	}

	return result
}

func (r *Registry) HistoryTotalRecordsForHour(monthsBack MonthsBack, day, hour int8) uint32 {
	if monthsBack >= 6 || day < 0 || day >= 31 || hour < 0 || hour >= 24 {
		return 0
	}

	slot, errGetHistorySlot := r.GetHistorySlot(monthsBack)
	if errGetHistorySlot != nil {
		return 0
	}

	return slot.GetMetric(day, hour).GetRecordsPerPeriod()
}
