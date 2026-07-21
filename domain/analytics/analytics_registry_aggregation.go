package analytics

func (r *Registry[T, TData]) aggregateHour(m *TMetric) AggregatedTopN {
	var result AggregatedTopN

	if (*m).GetRecordsPerPeriod() == 0 {
		return result
	}

	result.IPs = result.IPs.MergeActive(&m.TopIPs)
	result.ASN = result.ASN.MergeActive(&m.TopASN)
	result.Countries = result.Countries.MergeActive(&m.TopCountries)
	result.Cities = result.Cities.MergeActive(&m.TopCities)
	result.URL = result.URL.MergeActive(&m.TopURL)

	result.OS = result.OS.MergeActive(&m.TopOperatingSystems)
	result.Browsers = result.Browsers.MergeActive(&m.TopBrowsers)

	return result
}

func (r *Registry[T, TData]) CurrentMonthAggregateTopNForHour(day, hour int8) AggregatedTopN {
	if day < 0 || day >= 31 || hour < 0 || hour >= 24 {
		return AggregatedTopN{}
	}

	currentMonth := r.GetCurrentMonth()

	return r.aggregateHour((*currentMonth).GetMetric(day, hour))
}

func (r *Registry[T, TData]) PreviousMonthAggregateTopNForHour(day, hour int8) AggregatedTopN {
	if day < 0 || day >= 31 || hour < 0 || hour >= 24 {
		return AggregatedTopN{}
	}

	previousMonth := r.Ring.GetPreviousSlot()

	return r.aggregateHour((*previousMonth).GetMetric(day, hour))
}

func (r *Registry[T, TData]) PreviousMonthTotalRecords() uint32 {
	var result uint32

	previousMonth := r.Ring.GetPreviousSlot()

	for day := range int8(31) {
		for hour := range int8(24) {
			m := (*previousMonth).GetMetric(day, hour)

			result = result + (*m).GetRecordsPerPeriod()
		}
	}

	return result
}

func (r *Registry[T, TMetric]) PreviousMonthTotalRecordsForDay(day int8) uint32 {
	if day < 0 || day >= 31 {
		return 0
	}

	previousSlot := r.Ring.GetPreviousSlot()

	var result uint32

	r.ForEachMetric(
		previousSlot,
		func(d int8, h int8, m *TMetric) {
			if d == day {
				result = result + (*m).GetRecordsPerPeriod()
			}
		},
	)

	return result
}

func (r *Registry[T, TData]) CurrentMonthTotalRecords() uint32 {
	var result uint32

	currentSlot := r.GetCurrentMonth()

	for day := range int8(31) {
		for hour := range int8(24) {
			m := (*currentSlot).GetMetric(day, hour)

			result = result + (*m).GetRecordsPerPeriod()
		}
	}

	return result
}

func (r *Registry[T, TData]) CurrentMonthTotalRecordsForDay(day int8) uint32 {
	var result uint32

	if day < 0 || day >= 31 {
		return result
	}

	currentSlot := r.GetCurrentMonth()

	for hour := range int8(24) {
		m := (*currentSlot).GetMetric(day, hour)

		result = result + (*m).GetRecordsPerPeriod()
	}

	return result
}

func (r *Registry[T, TData]) mergeHourInto(m *TMetric, dst *AggregatedTopN) {
	if (*m).GetRecordsPerPeriod() == 0 {
		return
	}

	dst.IPs = dst.IPs.MergeActive(&m.TopIPs)
	dst.ASN = dst.ASN.MergeActive(&m.TopASN)
	dst.Countries = dst.Countries.MergeActive(&m.TopCountries)
	dst.Cities = dst.Cities.MergeActive(&m.TopCities)
	dst.URL = dst.URL.MergeActive(&m.TopURL)
	dst.OS = dst.OS.MergeActive(&m.TopOperatingSystems)
	dst.Browsers = dst.Browsers.MergeActive(&m.TopBrowsers)
}

func (r *Registry[T, Data]) PreviousMonthAggregateTopNForDay(day int8) AggregatedTopN {
	var result AggregatedTopN

	if day < 0 || day >= 31 {
		return result
	}

	previousSlot := r.Ring.GetPreviousSlot()

	for hour := range int8(24) {
		r.mergeHourInto((*previousSlot).GetMetric(day, hour), &result)
	}

	return result
}

func (r *Registry[T, D]) PreviousMonthAggregateTopN() AggregatedTopN {
	var result AggregatedTopN

	r.PreviousMonthForEach(
		func(_, _ int8, m *D) {
			r.mergeHourInto(m, &result)
		},
	)

	return result
}

func (r *Registry[T, Data]) CurrentMonthAggregateTopNForDay(day int8) AggregatedTopN {
	var result AggregatedTopN
	if day < 0 || day >= 31 {
		return result
	}

	currentSlot := r.GetCurrentMonth()

	for hour := range int8(24) {
		r.mergeHourInto((*currentSlot).GetMetric(day, hour), &result)
	}

	return result
}

func (r *Registry[T, Data]) CurrentMonthAggregateTopN() AggregatedTopN {
	var result AggregatedTopN

	r.CurrentMonthForEach(
		func(day int8, hour int8, m *MetricActive) {
			result.IPs = result.IPs.MergeActive(&m.TopIPs)
			result.ASN = result.ASN.MergeActive(&m.TopASN)
			result.Countries = result.Countries.MergeActive(&m.TopCountries)
			result.Cities = result.Cities.MergeActive(&m.TopCities)
			result.URL = result.URL.MergeActive(&m.TopURL)
			result.OS = result.OS.MergeActive(&m.TopOperatingSystems)
			result.Browsers = result.Browsers.MergeActive(&m.TopBrowsers)
		},
	)

	return result
}

func (r *Registry[T, TMetric]) HistoryAggregateTopNForDay(month, day int8) AggregatedTopN {
	var result AggregatedTopN

	if month < 0 || month >= 7 || day < 0 || day >= 31 {
		return result
	}

	slot := &r.Ring.slots[month]

	for hour := int8(0); hour < 24; hour++ {
		m := (*slot).GetMetric(day, hour)

		if (*m).GetRecordsPerPeriod() == 0 {
			continue
		}

		result.IPs = result.IPs.MergeActive(m.TopIPs)
		result.ASN = result.ASN.MergeArchived(m.TopASN)
		result.Countries = result.Countries.MergeArchived(m.TopCountries)
		result.Cities = result.Cities.MergeArchived(m.TopCities)
		result.URL = result.URL.MergeArchived(m.TopURL)
		result.OS = result.OS.MergeArchived(m.TopOperatingSystems)
		result.Browsers = result.Browsers.MergeArchived(m.TopBrowsers)
	}

	return result
}

func (r *Registry[T, TMetric]) HistoryAggregateTopNForMonth(month int8) AggregatedTopN {
	var result AggregatedTopN

	if month < 0 || month >= 7 {
		return result
	}

	slot := &r.Ring.slots[month]

	for day := int8(0); day < 31; day++ {
		for hour := int8(0); hour < 24; hour++ {
			m := (*slot).GetMetric(day, hour)

			if (*m).GetRecordsPerPeriod() == 0 {
				continue
			}

			result.IPs = result.IPs.MergeArchived(m.TopIPs)
			result.ASN = result.ASN.MergeArchived(m.TopASN)
			result.Countries = result.Countries.MergeArchived(m.TopCountries)
			result.Cities = result.Cities.MergeArchived(m.TopCities)
			result.URL = result.URL.MergeArchived(m.TopURL)
			result.OS = result.OS.MergeArchived(m.TopOperatingSystems)
			result.Browsers = result.Browsers.MergeArchived(m.TopBrowsers)
		}
	}

	return result
}

func (r *Registry[T, TMetric]) HistoryAggregateTopNForHour(month, day, hour int8) AggregatedTopN {
	var result AggregatedTopN

	if month < 0 || month >= 7 || day < 0 || day >= 31 || hour < 0 || hour >= 24 {
		return result
	}

	slot := &r.Ring.slots[month]

	m := (*slot).GetMetric(day, hour)

	if (*m).GetRecordsPerPeriod() == 0 {
		return result
	}

	result.IPs = result.IPs.MergeArchived(m.TopIPs)
	result.ASN = result.ASN.MergeArchived(m.TopASN)
	result.Countries = result.Countries.MergeArchived(m.TopCountries)
	result.Cities = result.Cities.MergeArchived(m.TopCities)
	result.URL = result.URL.MergeArchived(m.TopURL)
	result.OS = result.OS.MergeArchived(m.TopOperatingSystems)
	result.Browsers = result.Browsers.MergeArchived(m.TopBrowsers)

	return result
}

func (r *Registry[T, TMetric]) HistoryAggregateTopN() AggregatedTopN {
	var result AggregatedTopN

	for month := int8(0); month < 7; month++ {
		slot := &r.Ring.slots[month]

		for day := int8(0); day < 31; day++ {
			for hour := int8(0); hour < 24; hour++ {
				m := (*slot).GetMetric(day, hour)

				if (*m).GetRecordsPerPeriod() == 0 {
					continue
				}

				result.IPs = result.IPs.MergeArchived(m.TopIPs)
				result.ASN = result.ASN.MergeArchived(m.TopASN)
				result.Countries = result.Countries.MergeArchived(m.TopCountries)
				result.Cities = result.Cities.MergeArchived(m.TopCities)
				result.URL = result.URL.MergeArchived(m.TopURL)
				result.OS = result.OS.MergeArchived(m.TopOperatingSystems)
				result.Browsers = result.Browsers.MergeArchived(m.TopBrowsers)
			}
		}
	}

	return result
}

func (r *Registry[T, TMetric]) HistoryTotalRecords() uint32 {
	var result uint32

	for month := int8(0); month < 7; month++ {
		slot := &r.Ring.slots[month]

		for day := int8(0); day < 31; day++ {
			for hour := int8(0); hour < 24; hour++ {
				m := (*slot).GetMetric(day, hour)

				result = result + (*m).GetRecordsPerPeriod()
			}
		}
	}

	return result
}

func (r *Registry[T, TMetric]) HistoryTotalRecordsForMonth(month int8) uint32 {
	if month < 0 || month >= 7 {
		return 0
	}

	slot := &r.Ring.slots[month]

	var result uint32

	for day := int8(0); day < 31; day++ {
		for hour := int8(0); hour < 24; hour++ {
			m := (*slot).GetMetric(day, hour)

			result = result + (*m).GetRecordsPerPeriod()
		}
	}

	return result
}

func (r *Registry[T, TMetric]) HistoryTotalRecordsForDay(month, day int8) uint32 {
	if month < 0 || month >= 7 || day < 0 || day >= 31 {
		return 0
	}

	slot := &r.Ring.slots[month]

	var result uint32

	for hour := int8(0); hour < 24; hour++ {
		m := (*slot).GetMetric(day, hour)

		result = result + (*m).GetRecordsPerPeriod()
	}

	return result
}

func (r *Registry[T, TMetric]) HistoryTotalRecordsForHour(month, day, hour int8) uint32 {
	if month < 0 || month >= 7 || day < 0 || day >= 31 || hour < 0 || hour >= 24 {
		return 0
	}

	slot := &r.Ring.slots[month]

	m := (*slot).GetMetric(day, hour)

	return (*m).GetRecordsPerPeriod()
}
