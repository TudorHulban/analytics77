package site

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tudorhulban/analytics77/domain/analytics"
	"github.com/tudorhulban/analytics77/domain/offsets"
)

type SiteName string

// TODO: add methods to update the DST moments
type Site struct {
	Registry *analytics.Registry
	Name     SiteName `json:"name"`

	offsets.TimestampOffsets
}

func NewSite(rawSite map[string]any) (*Site, error) {
	valueReferenceMonth, exists := rawSite["referenceMonth"]
	if !exists {
		return nil,
			errors.New("missing referenceMonth")
	}
	floatReferenceMonth, couldCast := valueReferenceMonth.(float64)
	if !couldCast {
		return nil,
			errors.New("invalid referenceMonth type")
	}
	referenceMonth := int8(floatReferenceMonth)

	valueReferenceYear, exists := rawSite["referenceYear"]
	if !exists {
		return nil,
			errors.New("missing referenceYear")
	}
	floatReferenceYear, couldCast := valueReferenceYear.(float64)
	if !couldCast {
		return nil, errors.New("invalid referenceYear type")
	}
	referenceYear := int16(floatReferenceYear)

	data, errMarshal := json.Marshal(rawSite) // Marshal the individual site map back to JSON bytes
	if errMarshal != nil {
		return nil,
			fmt.Errorf(
				"failed to marshal site raw data: %w",
				errMarshal,
			)
	}

	var result Site

	if errUnMarshal := json.Unmarshal(data, &result); errUnMarshal != nil {
		return nil,
			fmt.Errorf(
				"failed to decode site configuration: %w",
				errUnMarshal,
			)
	}

	fmt.Println(referenceYear)

	result.Registry = analytics.NewRegistry(
		referenceMonth,
		result.TimestampDSTSpring,
		result.TimestampDSTWinter,
	)

	return &result,
		nil
}
