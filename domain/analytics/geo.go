package analytics

import (
	"fmt"
	"strings"
)

type GeoIP struct {
	ASN struct {
		AsNumber     string `json:"as_number"`
		Organization string `json:"organization"`
		Country      string `json:"country"`
	} `json:"asn"`

	Location struct {
		City        string `json:"city"`
		District    string `json:"district"`
		CountryCode string `json:"country_code3"`
		Postcode    string `json:"zipcode"`
		IsEU        bool   `json:"is_eu"`
	} `json:"location"`

	IsPrivate bool `json:"is_private"`
}

func (g GeoIP) String() string {
	if g.IsPrivate {
		return "GeoIP[Private Network]"
	}

	var sb strings.Builder

	// Location segment
	var locParts []string
	if g.Location.City != "" {
		locParts = append(locParts, g.Location.City)
	}

	if g.Location.District != "" {
		locParts = append(locParts, g.Location.District)
	}

	if g.Location.CountryCode != "" {
		locParts = append(locParts, g.Location.CountryCode)
	}

	sb.WriteString("GeoIP[")

	if len(locParts) > 0 {
		sb.WriteString(strings.Join(locParts, ", "))
	} else {
		sb.WriteString("Unknown Location")
	}

	if g.Location.Postcode != "" {
		fmt.Fprintf(&sb, " (%s)", g.Location.Postcode)
	}

	if g.Location.IsEU {
		sb.WriteString(" (EU)")
	}

	// ASN segment
	if g.ASN.AsNumber != "" || g.ASN.Organization != "" {
		sb.WriteString(" | ")

		if g.ASN.AsNumber != "" {
			sb.WriteString(g.ASN.AsNumber)
		}

		if g.ASN.Organization != "" {
			if g.ASN.AsNumber != "" {
				sb.WriteString(" - ")
			}

			sb.WriteString(g.ASN.Organization)
		}
	}

	sb.WriteString("]")

	return sb.String()
}
