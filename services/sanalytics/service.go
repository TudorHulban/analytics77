package sanalytics

import (
	"fmt"

	"github.com/tudorhulban/analytics77/domain/offsets"
	"github.com/tudorhulban/analytics77/infra/datacenter"
	"github.com/tudorhulban/analytics77/services/sgeo"
	"github.com/tudorhulban/analytics77/shared"
	"github.com/tudorhulban/hxhelpers/piers"
)

// ServiceAnalytics produces analytics for the datacenter.
//
// ServiceAnalytics can work with multiple sites.
type ServiceAnalytics struct {
	DC      *datacenter.DataCenter
	offsets *offsets.TimestampOffsets

	serviceGeo *sgeo.ServiceGeo
}

type PiersNewServiceAnalytics struct {
	ServiceGeo *sgeo.ServiceGeo
	Datacenter *datacenter.DataCenter
	Offsets    *offsets.TimestampOffsets
}

// ServiceAnalytics has a 1 to many relation with the datacenter due to the 1:1 relation with the geo service.
func NewServiceAnalytics(dependencies *PiersNewServiceAnalytics, siteName string) (*ServiceAnalytics, error) {
	if errValidate := piers.ValidateDependencies(dependencies); errValidate != nil {
		return nil,
			errValidate
	}

	if errDataCenter := dependencies.Datacenter.AddSite(siteName); errDataCenter != nil {
		return nil,
			errDataCenter
	}

	return &ServiceAnalytics{
			DC:         dependencies.Datacenter,
			serviceGeo: dependencies.ServiceGeo,

			offsets: dependencies.Offsets,
		},
		nil
}

// RecordEvents returns transformation and validation / processing errors.
func (s *ServiceAnalytics) RecordEvents(events shared.Requests) ([]error, []error) {
	errorsTransformation := make([]error, 0, len(events))
	validEvents := make([]*shared.ParamsAddEvent, 0, len(events))

	params := &shared.PiersAsParamsAddEvent{
		Offsets:    s.offsets,
		ServiceGeo: s.serviceGeo,
	}

	for ix, event := range events {
		param, errTransformation := event.AsParamsAddEvent(params)
		if errTransformation != nil {
			errorsTransformation = append(
				errorsTransformation,
				fmt.Errorf(
					"transformation error for event %d:%w",
					ix,
					errTransformation,
				),
			)

			continue
		}

		validEvents = append(validEvents, param)
	}

	if len(validEvents) == 0 {
		return errorsTransformation,
			nil
	}

	errorsProcess := s.DC.AddEvents(validEvents...)

	return errorsTransformation, errorsProcess
}
