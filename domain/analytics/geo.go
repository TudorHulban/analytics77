package analytics

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
