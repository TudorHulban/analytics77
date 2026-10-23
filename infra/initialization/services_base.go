package initialization

import (
	"os"

	"github.com/tudorhulban/analytics77/domain/offsets"
	"github.com/tudorhulban/analytics77/services/sanalytics"
	"github.com/tudorhulban/analytics77/services/sgeo"
	"github.com/tudorhulban/analytics77/services/slogging"
	"github.com/tudorhulban/analytics77/services/sstorage"
)

type ParamsServices struct {
	ServiceLogging *slogging.ServiceLogging

	APIKeyGeolocation string
	Offsets           offsets.TimestampOffsets
}

func Services(params *ParamsServices) (*sanalytics.ServiceAnalytics, error) {
	params.
		ServiceLogging.
		Logger.Print("Initializing Analytics Application...")

	serviceStorage, errCrServiceStorage := sstorage.NewServiceStorage(".")
	if errCrServiceStorage != nil {
		params.
			ServiceLogging.
			Logger.Printf(
			"service geo creation: %s\n",
			errCrServiceStorage.Error(),
		)

		os.Exit(10)
	}

	serviceGeo, errCrServiceGeo := sgeo.NewServiceGeo(
		&sgeo.ParamsNewServiceGeo{
			APIKeyGeolocation: params.APIKeyGeolocation,
		},
		serviceStorage,
	)
	if errCrServiceGeo != nil {
		params.
			ServiceLogging.
			Logger.Printf(
			"service geo creation: %s",
			errCrServiceGeo.Error(),
		)

		os.Exit(11)
	}

	result, errCrServiceAnalytics := sanalytics.NewServiceAnalytics(
		&sanalytics.PiersNewServiceAnalytics{
			ServiceGeo: serviceGeo,
			Offsets:    &params.Offsets,
		},
		"some-ste", // TODO: change
	)
	if errCrServiceAnalytics != nil {
		params.
			ServiceLogging.
			Logger.Printf(
			"service analytics creation: %s",
			errCrServiceAnalytics.Error(),
		)

		os.Exit(11)
	}

	return result,
		nil
}
