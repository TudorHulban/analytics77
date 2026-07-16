package analytics

func (*Registry) aggregateHour(m *MetricActive) AggregatedTopN {
	var result AggregatedTopN

	if m.RecordsPerPeriod.Load() == 0 {
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

func (r *Registry) CurrentMonthAggregateTopNForHour(day, hour int8) AggregatedTopN {
	if day < 0 || day >= 31 || hour < 0 || hour >= 24 {
		return AggregatedTopN{}
	}

	return r.aggregateHour(&r.GetCurrentMonth()[day][hour])
}

func (r *Registry) PreviousMonthAggregateTopNForHour(day, hour int8) AggregatedTopN {
	if day < 0 || day >= 31 || hour < 0 || hour >= 24 {
		return AggregatedTopN{}
	}

	return r.aggregateHour(&r.GetPreviousMonth()[day][hour])
}

func (r *Registry) PreviousMonthTotalRecords() uint32 {
	var result uint32

	for day := range int8(31) {
		for hour := range int8(24) {
			m := &r.GetPreviousMonth()[day][hour]

			result = result + m.RecordsPerPeriod.Load()
		}
	}

	return result
}

func (r *Registry) PreviousMonthTotalRecordsForDay(day int8) uint32 {
	var result uint32

	if day < 0 || day >= 31 {
		return result
	}

	for hour := range int8(24) {
		m := &r.GetPreviousMonth()[day][hour]

		result = result + m.RecordsPerPeriod.Load()
	}

	return result
}

func (r *Registry) CurrentMonthTotalRecords() uint32 {
	var result uint32

	for day := range int8(31) {
		for hour := range int8(24) {
			m := &r.GetCurrentMonth()[day][hour]

			result = result + m.RecordsPerPeriod.Load()
		}
	}

	return result
}

func (r *Registry) CurrentMonthTotalRecordsForDay(day int8) uint32 {
	var result uint32

	if day < 0 || day >= 31 {
		return result
	}

	for hour := range int8(24) {
		m := &r.GetCurrentMonth()[day][hour]

		result = result + m.RecordsPerPeriod.Load()
	}

	return result
}

func (*Registry) mergeHourInto(m *MetricActive, dst *AggregatedTopN) {
	if m.RecordsPerPeriod.Load() == 0 {
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

func (r *Registry) PreviousMonthAggregateTopNForDay(day int8) AggregatedTopN {
	var result AggregatedTopN

	if day < 0 || day >= 31 {
		return result
	}

	for hour := range 24 {
		r.mergeHourInto(&r.GetPreviousMonth()[day][hour], &result)
	}

	return result
}

func (r *Registry) PreviousMonthAggregateTopN() AggregatedTopN {
	var result AggregatedTopN

	r.PreviousMonthForEach(
		func(_, _ int8, m *MetricActive) {
			r.mergeHourInto(m, &result)
		},
	)

	return result
}

func (r *Registry) CurrentMonthAggregateTopNForDay(day int8) AggregatedTopN {
	var result AggregatedTopN
	if day < 0 || day >= 31 {
		return result
	}

	for hour := range 24 {
		r.mergeHourInto(&r.GetCurrentMonth()[day][hour], &result)
	}

	return result
}

func (r *Registry) CurrentMonthAggregateTopN() AggregatedTopN {
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

func (r *Registry) HistoryAggregateTopNForDay(month, day int8) AggregatedTopN {
	var result AggregatedTopN

	if month < 0 || month >= 7 || day < 0 || day >= 31 {
		return result
	}

	for hour := range 24 {
		m := &r.History[month][day][hour]

		if m.RecordsPerPeriod == 0 {
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

	return result
}

func (r *Registry) HistoryAggregateTopNForMonth(month int8) AggregatedTopN {
	var result AggregatedTopN

	if month < 0 || month >= 7 {
		return result
	}

	for day := range 31 {
		for hour := range 24 {
			m := &r.History[month][day][hour]

			if m.RecordsPerPeriod == 0 {
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

func (r *Registry) HistoryAggregateTopNForHour(month, day, hour int8) AggregatedTopN {
	var result AggregatedTopN

	if month < 0 || month >= 7 || day < 0 || day >= 31 || hour < 0 || hour >= 24 {
		return result
	}

	m := &r.History[month][day][hour]

	if m.RecordsPerPeriod == 0 {
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

func (r *Registry) HistoryAggregateTopN() AggregatedTopN {
	var result AggregatedTopN

	for month := range 7 {
		for day := range 31 {
			for hour := range 24 {
				m := &r.History[month][day][hour]

				if m.RecordsPerPeriod == 0 {
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

func (r *Registry) HistoryTotalRecords() uint32 {
	var result uint32

	for month := range 7 {
		for day := range 31 {
			for hour := range 24 {
				result = result + r.History[month][day][hour].RecordsPerPeriod
			}
		}
	}

	return result
}

func (r *Registry) HistoryTotalRecordsForMonth(month int8) uint32 {
	var result uint32

	if month < 0 || month >= 7 {
		return result
	}

	for day := range 31 {
		for hour := range 24 {
			result = result + r.History[month][day][hour].RecordsPerPeriod
		}
	}

	return result
}

func (r *Registry) HistoryTotalRecordsForDay(month, day int8) uint32 {
	var result uint32

	if month < 0 || month >= 7 || day < 0 || day >= 31 {
		return result
	}

	for hour := range 24 {
		result = result + r.History[month][day][hour].RecordsPerPeriod
	}

	return result
}

func (r *Registry) HistoryTotalRecordsForHour(month, day, hour int8) uint32 {
	if month < 0 || month >= 7 || day < 0 || day >= 31 || hour < 0 || hour >= 24 {
		return 0
	}

	return r.History[month][day][hour].RecordsPerPeriod
}
