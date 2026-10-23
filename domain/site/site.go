package site

import (
	"github.com/tudorhulban/analytics77/domain/analytics"
	"github.com/tudorhulban/analytics77/domain/offsets"
)

type SiteName string

// TODO: add methods to update the DST moments
type Site struct {
	Registry *analytics.Registry
	Name     SiteName

	offsets.TimestampOffsets
}
