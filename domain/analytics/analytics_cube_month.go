package analytics

type (
	MonthActive   [31]DayActive
	MonthArchived [31]DayArchived
)

func (m *MonthActive) IsZero() bool {
	for i := range m {
		if !m[i].IsZero() {
			return false // Found data in one of the days
		}
	}

	return true // All 31 days are completely zero
}

func (m *MonthActive) DeepCopyInto(dst *MonthActive) {
	for d := range 31 {
		for h := range 24 {
			m[d][h].TopIPs.DeepCopyInto(&dst[d][h].TopIPs)
			m[d][h].TopASN.DeepCopyInto(&dst[d][h].TopASN)
			m[d][h].TopCountries.DeepCopyInto(&dst[d][h].TopCountries)
			m[d][h].TopCities.DeepCopyInto(&dst[d][h].TopCities)
			m[d][h].TopURL.DeepCopyInto(&dst[d][h].TopURL)
			m[d][h].TopOperatingSystems.DeepCopyInto(&dst[d][h].TopOperatingSystems)
			m[d][h].TopBrowsers.DeepCopyInto(&dst[d][h].TopBrowsers)

			dst[d][h].RecordsPerPeriod.Store(m[d][h].RecordsPerPeriod.Load())
		}
	}
}
