package site

import (
	"errors"
	"fmt"

	"github.com/tudorhulban/analytics77/domain/analytics"
	"github.com/tudorhulban/analytics77/domain/offsets"
)

type SiteName string

// TODO: add methods to update the DST moments
type Site struct {
	Registry *analytics.Registry `json:"-"`
	Name     SiteName            `json:"name"`

	offsets.TimestampOffsets `json:"-"`
}

func NewSite(rawSite map[string]any) (*Site, error) {
	valueReferenceMonth, exists := rawSite["referenceMonth"]
	if !exists {
		return nil,
			errors.New("missing referenceMonth")
	}

	floatReferenceMonth, couldCastReferenceMonth := valueReferenceMonth.(float64)
	if !couldCastReferenceMonth {
		return nil,
			errors.New("invalid referenceMonth type")
	}

	referenceMonth := int8(floatReferenceMonth)

	valueReferenceYear, exists := rawSite["referenceYear"]
	if !exists {
		return nil,
			errors.New("missing referenceYear")
	}

	floatReferenceYear, couldCastReferenceYear := valueReferenceYear.(float64)
	if !couldCastReferenceYear {
		return nil,
			fmt.Errorf(
				"invalid referenceYear type %f",
				floatReferenceYear,
			)
	}

	_ = int16(floatReferenceYear) // referenceYear

	valueSiteName, exists := rawSite["name"]
	if !exists {
		return nil,
			errors.New("missing site name")
	}

	siteName, couldSiteName := valueSiteName.(string)
	if !couldSiteName {
		return nil,
			fmt.Errorf(
				"invalid site name %v",
				valueSiteName,
			)
	}

	valueOffsetUTCHours, exists := rawSite["offsetUTCHours"]
	if !exists {
		return nil,
			errors.New("missing offsetUTCHours")
	}

	offsetUTCHours, couldCastUTCHours := valueOffsetUTCHours.(float64)
	if !couldCastUTCHours {
		return nil,
			fmt.Errorf(
				"invalid offsetUTCHours %v",
				valueOffsetUTCHours,
			)
	}

	valueDSTSpring, exists := rawSite["timestampDSTSpring"]
	if !exists {
		return nil,
			errors.New("missing timestampDSTSpring")
	}

	timestampDSTSpring, couldCastDSTSpring := valueDSTSpring.(float64)
	if !couldCastDSTSpring {
		return nil,
			fmt.Errorf(
				"invalid timestampDSTSpring %v",
				valueDSTSpring,
			)
	}

	valueDSTWinter, exists := rawSite["timestampDSTWinter"]
	if !exists {
		return nil,
			errors.New("missing timestampDSTWinter")
	}

	timestampDSTWinter, couldCastDSTWinter := valueDSTWinter.(float64)
	if !couldCastDSTWinter {
		return nil,
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

	result.Registry = analytics.NewRegistry(
		referenceMonth,
		result.TimestampDSTSpring,
		result.TimestampDSTWinter,
	)

	return &result,
		nil
}
