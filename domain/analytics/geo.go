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

	var partsLocation []string
	if g.Location.City != "" {
		partsLocation = append(partsLocation, g.Location.City)
	}

	if g.Location.District != "" {
		partsLocation = append(partsLocation, g.Location.District)
	}

	if g.Location.CountryCode != "" {
		partsLocation = append(partsLocation, g.Location.CountryCode)
	}

	sb.WriteString("GeoIP[")

	if len(partsLocation) > 0 {
		sb.WriteString(strings.Join(partsLocation, ", "))
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
