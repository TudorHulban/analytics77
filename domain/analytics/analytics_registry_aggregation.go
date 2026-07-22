package analytics

func (r *Registry) aggregateHour(fromMetric *MetricActive) *AggregatedTopN {
	var result AggregatedTopN

	if (*fromMetric).GetRecordsPerPeriod() == 0 {
		return &result
	}

	result.IPs = *(*fromMetric).GetTopIPs()
	result.ASN = *(*fromMetric).GetTopASNs()
	result.Countries = *(*fromMetric).GetTopCountries()
	result.Cities = *(*fromMetric).GetTopCities()
	result.URL = *(*fromMetric).GetTopURLs()

	result.OS = *(*fromMetric).GetTopOperatingSystems()
	result.Browsers = *(*fromMetric).GetTopBrowsers()

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
			m := (*previousMonth).GetMetric(day, hour)

			result = result + (*m).GetRecordsPerPeriod()
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
			m := (*currentSlot).GetMetric(day, hour)

			result = result + (*m).GetRecordsPerPeriod()
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

func (r *Registry) mergeHourInto(fromMetric *MetricActive, dst *AggregatedTopN) {
	if (*fromMetric).GetRecordsPerPeriod() == 0 {
		return
	}

	dst.IPs = *(*fromMetric).GetTopIPs()
	dst.ASN = *(*fromMetric).GetTopASNs()
	dst.Countries = *(*fromMetric).GetTopCountries()
	dst.Cities = *(*fromMetric).GetTopCities()
	dst.URL = *(*fromMetric).GetTopURLs()

	dst.OS = *(*fromMetric).GetTopOperatingSystems()
	dst.Browsers = *(*fromMetric).GetTopBrowsers()
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
		func(day int8, hour int8, fromMetric *MetricActive) {
			result.IPs = *(*fromMetric).GetTopIPs()
			result.ASN = *(*fromMetric).GetTopASNs()
			result.Countries = *(*fromMetric).GetTopCountries()
			result.Cities = *(*fromMetric).GetTopCities()
			result.URL = *(*fromMetric).GetTopURLs()

			result.OS = *(*fromMetric).GetTopOperatingSystems()
			result.Browsers = *(*fromMetric).GetTopBrowsers()
		},
	)

	return &result
}

func (r *Registry) HistoryAggregateTopNForDay(month, day int8) (*AggregatedTopN, error) {
	if month < 0 || month >= 7 || day < 0 || day >= 31 {
		return nil,
			ErrInvalidInput
	}

	var result AggregatedTopN

	slot := &r.Slots[month]

	for hour := int8(0); hour < 24; hour++ {
		fromMetric := (*slot).GetMetric(day, hour)

		if (*fromMetric).GetRecordsPerPeriod() == 0 {
			continue
		}

		result.IPs = *(*fromMetric).GetTopIPs()
		result.ASN = *(*fromMetric).GetTopASNs()
		result.Countries = *(*fromMetric).GetTopCountries()
		result.Cities = *(*fromMetric).GetTopCities()
		result.URL = *(*fromMetric).GetTopURLs()

		result.OS = *(*fromMetric).GetTopOperatingSystems()
		result.Browsers = *(*fromMetric).GetTopBrowsers()
	}

	return &result, nil
}

func (r *Registry) HistoryAggregateTopNForMonth(month int8) (*AggregatedTopN, error) {
	if month < 0 || month >= 7 {
		return nil,
			ErrInvalidInput
	}

	var result AggregatedTopN

	slot := &r.Slots[month]

	for day := int8(0); day < 31; day++ {
		for hour := int8(0); hour < 24; hour++ {
			fromMetric := (*slot).GetMetric(day, hour)

			if (*fromMetric).GetRecordsPerPeriod() == 0 {
				continue
			}

			result.IPs = *(*fromMetric).GetTopIPs()
			result.ASN = *(*fromMetric).GetTopASNs()
			result.Countries = *(*fromMetric).GetTopCountries()
			result.Cities = *(*fromMetric).GetTopCities()
			result.URL = *(*fromMetric).GetTopURLs()

			result.OS = *(*fromMetric).GetTopOperatingSystems()
			result.Browsers = *(*fromMetric).GetTopBrowsers()
		}
	}

	return &result, nil
}

func (r *Registry) HistoryAggregateTopNForHour(month, day, hour int8) (*AggregatedTopN, error) {
	if month < 0 || month >= 7 || day < 0 || day >= 31 || hour < 0 || hour >= 24 {
		return nil,
			ErrInvalidInput
	}

	var result AggregatedTopN

	slot := &r.Slots[month]

	fromMetric := (*slot).GetMetric(day, hour)

	if (*fromMetric).GetRecordsPerPeriod() == 0 {
		return &result, nil
	}

	result.IPs = *(*fromMetric).GetTopIPs()
	result.ASN = *(*fromMetric).GetTopASNs()
	result.Countries = *(*fromMetric).GetTopCountries()
	result.Cities = *(*fromMetric).GetTopCities()
	result.URL = *(*fromMetric).GetTopURLs()

	result.OS = *(*fromMetric).GetTopOperatingSystems()
	result.Browsers = *(*fromMetric).GetTopBrowsers()

	return &result, nil
}

func (r *Registry) HistoryAggregateTopN() *AggregatedTopN {
	var result AggregatedTopN

	for month := int8(0); month < 7; month++ {
		slot := &r.Slots[month]

		for day := int8(0); day < 31; day++ {
			for hour := int8(0); hour < 24; hour++ {
				fromMetric := (*slot).GetMetric(day, hour)

				if (*fromMetric).GetRecordsPerPeriod() == 0 {
					continue
				}

				result.IPs = *(*fromMetric).GetTopIPs()
				result.ASN = *(*fromMetric).GetTopASNs()
				result.Countries = *(*fromMetric).GetTopCountries()
				result.Cities = *(*fromMetric).GetTopCities()
				result.URL = *(*fromMetric).GetTopURLs()

				result.OS = *(*fromMetric).GetTopOperatingSystems()
				result.Browsers = *(*fromMetric).GetTopBrowsers()
			}
		}
	}

	return &result
}

func (r *Registry) HistoryTotalRecords() uint32 {
	var result uint32

	for month := int8(0); month < 7; month++ {
		slot := &r.Slots[month]

		for day := int8(0); day < 31; day++ {
			for hour := int8(0); hour < 24; hour++ {
				m := (*slot).GetMetric(day, hour)

				result = result + (*m).GetRecordsPerPeriod()
			}
		}
	}

	return result
}

func (r *Registry) HistoryTotalRecordsForMonth(month int8) uint32 {
	if month < 0 || month >= 7 {
		return 0
	}

	slot := r.Slots[month]

	var result uint32

	for day := int8(0); day < 31; day++ {
		for hour := int8(0); hour < 24; hour++ {
			m := slot.GetMetric(day, hour)

			result = result + (*m).GetRecordsPerPeriod()
		}
	}

	return result
}

func (r *Registry) HistoryTotalRecordsForDay(month, day int8) uint32 {
	if month < 0 || month >= 7 || day < 0 || day >= 31 {
		return 0
	}

	slot := r.Slots[month]

	var result uint32

	for hour := int8(0); hour < 24; hour++ {
		m := slot.GetMetric(day, hour)

		result = result + (*m).GetRecordsPerPeriod()
	}

	return result
}

func (r *Registry) HistoryTotalRecordsForHour(month, day, hour int8) uint32 {
	if month < 0 || month >= 7 || day < 0 || day >= 31 || hour < 0 || hour >= 24 {
		return 0
	}

	slot := r.Slots[month]

	m := slot.GetMetric(day, hour)

	return (*m).GetRecordsPerPeriod()
}
