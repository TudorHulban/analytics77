package shared

import (
	"errors"

	"github.com/tudorhulban/analytics77/domain/analytics"
)

type DayMonth uint8 // 1 - 31

func (d DayMonth) IsValid() bool {
	return d >= 1 && d <= 31
}

type HourDay uint8 // 0 - 23

func (h HourDay) IsValid() bool {
	return h <= 23
}

type ParamsAddEvent struct {
	SiteKey         string
	Country         string
	City            string
	IP              string
	ASNOrganization string

	TimestampUNIX  int64
	OffsetUTCHours int64

	Browser         analytics.Browser
	OperatingSystem analytics.OS

	IsPrivateIP bool
}

func (e *ParamsAddEvent) Validate() []error {
	var errs []error

	if len(e.SiteKey) == 0 {
		errs = append(
			errs,
			errors.New("site key cannot be empty"),
		)
	}

	if !e.IsPrivateIP {
		if len(e.Country) == 0 {
			errs = append(
				errs,
				errors.New("country cannot be empty"),
			)
		}

		if len(e.City) == 0 {
			errs = append(
				errs,
				errors.New("city cannot be empty"),
			)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}
