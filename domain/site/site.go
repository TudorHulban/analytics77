package site

import (
	"errors"
	"fmt"

	"github.com/tudorhulban/analytics77/domain/offsets"
)

type SiteName string

// Site is based on primitives so no need to pass by reference externally.
type Site struct {
	Name         SiteName `json:"name"`
	SiteLocation string   // IANA Location.
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

	valueOffsetUTCHours, exists := rawSite["offsetUTCHours"]
	if !exists {
		return nil,
			0,
			errors.New("missing offsetUTCHours")
	}

	offsetUTCHours, couldCastUTCHours := valueOffsetUTCHours.(float64)
	if !couldCastUTCHours {
		return nil,
			0,
			fmt.Errorf(
				"invalid offsetUTCHours %v",
				valueOffsetUTCHours,
			)
	}

	valueDSTSpring, exists := rawSite["timestampDSTSpring"]
	if !exists {
		return nil,
			0,
			errors.New("missing timestampDSTSpring")
	}

	timestampDSTSpring, couldCastDSTSpring := valueDSTSpring.(float64)
	if !couldCastDSTSpring {
		return nil,
			0,
			fmt.Errorf(
				"invalid timestampDSTSpring %v",
				valueDSTSpring,
			)
	}

	valueDSTWinter, exists := rawSite["timestampDSTWinter"]
	if !exists {
		return nil,
			0,
			errors.New("missing timestampDSTWinter")
	}

	timestampDSTWinter, couldCastDSTWinter := valueDSTWinter.(float64)
	if !couldCastDSTWinter {
		return nil,
			0,
			fmt.Errorf(
				"invalid timestampDSTSpring %v",
				valueDSTSpring,
			)
	}

	result := Site{
		Name: SiteName(siteName),

		TimestampOffsets: offsets.TimestampOffsets{
			OffsetUTCHours:     int64(offsetUTCHours),
			TimestampDSTSpring: int64(timestampDSTSpring),
			TimestampDSTWinter: int64(timestampDSTWinter),
		},
	}

	return &result,
		referenceMonth,
		nil
}
