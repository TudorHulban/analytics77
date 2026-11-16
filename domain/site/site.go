package site

import (
	"errors"
	"fmt"
)

type SiteName string

// Site is based on primitives so no need to pass by reference externally.
type Site struct {
	Name         SiteName `json:"name"`
	SiteLocation string   `json:"location_iana"` // IANA Location.
}

func NewSite(rawSite map[string]any) (*Site, int8, error) {
	valueReferenceMonth, exists := rawSite["referenceMonth"]
	if !exists {
		return nil,
			0,
			errors.New("missing referenceMonth")
	}

	floatReferenceMonth, couldCastReferenceMonth := valueReferenceMonth.(float64)
	if !couldCastReferenceMonth {
		return nil,
			0,
			errors.New("invalid referenceMonth type")
	}

	referenceMonth := int8(floatReferenceMonth)

	valueReferenceYear, exists := rawSite["referenceYear"]
	if !exists {
		return nil,
			0,
			errors.New("missing referenceYear")
	}

	floatReferenceYear, couldCastReferenceYear := valueReferenceYear.(float64)
	if !couldCastReferenceYear {
		return nil,
			0,
			fmt.Errorf(
				"invalid referenceYear type %f",
				floatReferenceYear,
			)
	}

	_ = int16(floatReferenceYear) // referenceYear

	valueSiteName, exists := rawSite["name"]
	if !exists {
		return nil,
			0,
			errors.New("missing site name")
	}

	siteName, couldSiteName := valueSiteName.(string)
	if !couldSiteName {
		return nil,
			0,
			fmt.Errorf(
				"invalid site name %v",
				valueSiteName,
			)
	}

	valueIANALocation, exists := rawSite["location_iana"]
	if !exists {
		return nil,
			0,
			errors.New("missing IANA location")
	}

	locationIANA, couldCastLocationIANA := valueIANALocation.(string)
	if !couldCastLocationIANA {
		return nil,
			0,
			fmt.Errorf(
				"invalid site name %v",
				valueSiteName,
			)
	}

	result := Site{
		Name:         SiteName(siteName),
		SiteLocation: locationIANA,
	}

	return &result,
		referenceMonth,
		nil
}
