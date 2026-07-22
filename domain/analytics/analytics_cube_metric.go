package analytics

import "sync/atomic"

type MetricActive struct {
	TopIPs              MetaActive[string]
	TopASN              MetaActive[string]
	TopCountries        MetaActive[string]
	TopCities           MetaActive[string]
	TopURL              MetaActive[string] // TODO: TopURLs
	TopOperatingSystems MetaActive[OS]
	TopBrowsers         MetaActive[Browser]

	RecordsPerPeriod atomic.Uint32
	readOnly         atomic.Bool
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
	return &m.TopCountries
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
		m.TopURL.occupied.Load() != 0 ||
		m.TopOperatingSystems.occupied.Load() != 0 ||
		m.TopBrowsers.occupied.Load() != 0 {
		return false
	}

	return true
}

func (m *MetricActive) DeepCopyInto(dst *MetricActive) {
	m.readOnly.Store(true)

	m.TopIPs.DeepCopyInto(&dst.TopIPs)
	m.TopASN.DeepCopyInto(&dst.TopASN)
	m.TopCountries.DeepCopyInto(&dst.TopCountries)
	m.TopCities.DeepCopyInto(&dst.TopCities)
	m.TopURL.DeepCopyInto(&dst.TopURL)
	m.TopOperatingSystems.DeepCopyInto(&dst.TopOperatingSystems)
	m.TopBrowsers.DeepCopyInto(&dst.TopBrowsers)

	dst.RecordsPerPeriod.Store(m.RecordsPerPeriod.Load())

	m.readOnly.Store(false)
}
