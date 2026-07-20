package analytics

import "sync/atomic"

type MetricActive struct {
	TopIPs              MetaActive[string]
	TopASN              MetaActive[string]
	TopCountries        MetaActive[string]
	TopCities           MetaActive[string]
	TopURL              MetaActive[string]
	TopOperatingSystems MetaActive[OS]
	TopBrowsers         MetaActive[Browser]

	RecordsPerPeriod atomic.Uint32
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

func (m *MetricActive) AsArchived() MetricArchived {
	return MetricArchived{
		RecordsPerPeriod:    m.RecordsPerPeriod.Load(),
		TopIPs:              m.TopIPs.AsMetaArchive(),
		TopASN:              m.TopASN.AsMetaArchive(),
		TopCountries:        m.TopCountries.AsMetaArchive(),
		TopCities:           m.TopCities.AsMetaArchive(),
		TopURL:              m.TopURL.AsMetaArchive(),
		TopOperatingSystems: m.TopOperatingSystems.AsMetaArchive(),
		TopBrowsers:         m.TopBrowsers.AsMetaArchive(),
	}
}

func (m *MetricActive) DeepCopyInto(dst *MetricActive) {
	m.TopIPs.DeepCopyInto(&dst.TopIPs)
	m.TopASN.DeepCopyInto(&dst.TopASN)
	m.TopCountries.DeepCopyInto(&dst.TopCountries)
	m.TopCities.DeepCopyInto(&dst.TopCities)
	m.TopURL.DeepCopyInto(&dst.TopURL)
	m.TopOperatingSystems.DeepCopyInto(&dst.TopOperatingSystems)
	m.TopBrowsers.DeepCopyInto(&dst.TopBrowsers)

	dst.RecordsPerPeriod.Store(m.RecordsPerPeriod.Load())
}

type MetricArchived struct {
	TopIPs              MetaArchived[string]
	TopASN              MetaArchived[string]
	TopCountries        MetaArchived[string]
	TopCities           MetaArchived[string]
	TopURL              MetaArchived[string]
	TopOperatingSystems MetaArchived[OS]
	TopBrowsers         MetaArchived[Browser]

	RecordsPerPeriod uint32
}
