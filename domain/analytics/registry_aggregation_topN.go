package analytics

import "strings"

type AggregatedTopN struct {
	IPs       MetaActive[string]
	ASN       MetaActive[string]
	Countries MetaActive[string]
	Cities    MetaActive[string]
	URL       MetaActive[string]
	OS        MetaActive[OS]
	Browsers  MetaActive[Browser]
}

func (a *AggregatedTopN) IsZero() bool {
	if a == nil {
		return true
	}

	return a.IPs.IsZero() &&
		a.ASN.IsZero() &&
		a.Countries.IsZero() &&
		a.Cities.IsZero() &&
		a.URL.IsZero() &&
		a.OS.IsZero() &&
		a.Browsers.IsZero()
}

func (a *AggregatedTopN) String() string {
	var b strings.Builder

	b.WriteString("IPs:\n")
	b.WriteString(a.IPs.String())
	// b.WriteByte('\n')

	b.WriteString("ASN:\n")
	b.WriteString(a.ASN.String())
	// b.WriteByte('\n')

	b.WriteString("Countries:\n")
	b.WriteString(a.Countries.String())
	// b.WriteByte('\n')

	b.WriteString("Cities:\n")
	b.WriteString(a.Cities.String())
	// b.WriteByte('\n')

	b.WriteString("URL:\n")
	b.WriteString(a.URL.String())
	// b.WriteByte('\n')

	b.WriteString("OS:\n")
	b.WriteString(a.OS.String())
	// b.WriteByte('\n')

	b.WriteString("Browsers:\n")
	b.WriteString(a.Browsers.String())
	// b.WriteByte('\n')

	return b.String()
}
