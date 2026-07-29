package analytics

import "sync/atomic"

type MetricActive struct {
	TopIPs              MetaActive[string]
	TopASN              MetaActive[string]
	TopCountries        MetaActive[string]
	TopCities           MetaActive[string]
	TopURLs             MetaActive[string]
	TopOperatingSystems MetaActive[OS]
	TopBrowsers         MetaActive[Browser]

	RecordsPerPeriod atomic.Uint32
}

func (m *MetricActive) GetRecordsPerPeriod() uint32 {
	return m.RecordsPerPeriod.Load()
}

func (m *MetricActive) GetTopIPs() *MetaActive[string] {
	return &m.TopIPs
}

func (m *MetricActive) GetTopASNs() *MetaActive[string] {
	return &m.TopASN
}

func (m *MetricActive) GetTopCountries() *MetaActive[string] {
	return &m.TopCountries
}

func (m *MetricActive) GetTopCities() *MetaActive[string] {
	return &m.TopCities
}

func (m *MetricActive) GetTopURLs() *MetaActive[string] {
	return &m.TopURLs
}

func (m *MetricActive) GetTopOperatingSystems() *MetaActive[OS] {
	return &m.TopOperatingSystems
}

func (m *MetricActive) GetTopBrowsers() *MetaActive[Browser] {
	return &m.TopBrowsers
}

func (m *MetricActive) IsZero() bool {
	// 1. Check the atomic counter first (usually fastest to rule out)
	if m.RecordsPerPeriod.Load() != 0 {
		return false
	}

	// 2. Check each MetaActive field via their occupied bits
	if m.TopIPs.occupied.Load() != 0 ||
		m.TopASN.occupied.Load() != 0 ||
		m.TopCountries.occupied.Load() != 0 ||
		m.TopCities.occupied.Load() != 0 ||
		m.TopURLs.occupied.Load() != 0 ||
		m.TopOperatingSystems.occupied.Load() != 0 ||
		m.TopBrowsers.occupied.Load() != 0 {
		return false
	}

	return true
}

func (m *MetricActive) deepCopyInto(dst *MetricActive) {
	m.TopIPs.deepCopyInto(&dst.TopIPs)
	m.TopASN.deepCopyInto(&dst.TopASN)
	m.TopCountries.deepCopyInto(&dst.TopCountries)
	m.TopCities.deepCopyInto(&dst.TopCities)
	m.TopURLs.deepCopyInto(&dst.TopURLs)
	m.TopOperatingSystems.deepCopyInto(&dst.TopOperatingSystems)
	m.TopBrowsers.deepCopyInto(&dst.TopBrowsers)

	dst.RecordsPerPeriod.Store(m.RecordsPerPeriod.Load())
}
