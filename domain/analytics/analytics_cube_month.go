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
			m[d][h].DeepCopyInto(&dst[d][h])
		}
	}
}
