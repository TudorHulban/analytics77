package sanalytics

import (
	"fmt"

	"github.com/tudorhulban/analytics77/helpers"
	"github.com/tudorhulban/analytics77/infra/datacenter"
	"github.com/tudorhulban/analytics77/services/sgeo"
	"github.com/tudorhulban/analytics77/shared"
	"github.com/tudorhulban/hxhelpers/piers"
)

// TODO: add methods to update the DST moments
type ServiceAnalytics struct {
	DC      *datacenter.DataCenter
	offsets *helpers.TimestampOffsets

	serviceGeo *sgeo.ServiceGeo
}

type PiersNewServiceAnalytics struct {
	ServiceGeo *sgeo.ServiceGeo
}

func NewServiceAnalytics(dependencies *PiersNewServiceAnalytics, offsets *helpers.TimestampOffsets) (*ServiceAnalytics, error) {
	if errValidate := piers.ValidateDependencies(dependencies); errValidate != nil {
		return nil,
			errValidate
	}

	return &ServiceAnalytics{
			DC:         datacenter.NewDataCenter(),
			serviceGeo: dependencies.ServiceGeo,

			offsets: offsets,
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
