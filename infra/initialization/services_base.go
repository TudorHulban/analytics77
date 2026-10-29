package initialization

import (
	"fmt"
	"os"

	"github.com/tudorhulban/analytics77/domain/offsets"
	"github.com/tudorhulban/analytics77/services/sanalytics"
	"github.com/tudorhulban/analytics77/services/sgeo"
	"github.com/tudorhulban/analytics77/services/sstorage"
	"github.com/tudorhulban/hxerrors"
)

type ParamsServices struct {
	APIKeyGeolocation string
	Offsets           offsets.TimestampOffsets
}

func NewApplicationServices(params *ParamsServices) (*sanalytics.ServiceAnalytics, error) {
	fmt.Println("Initializing Analytics Application...")

	serviceStorage, errCrServiceStorage := sstorage.NewServiceStorage(".")
	if errCrServiceStorage != nil {
		fmt.Printf(
			"creation service storage: %s\n",
			errCrServiceStorage.Error(),
		)

		os.Exit(hxerrors.OSExitForStorageIssues)
	}

	serviceGeo, errCrServiceGeo := sgeo.NewServiceGeo(
		&sgeo.ParamsNewServiceGeo{
			APIKeyGeolocation: params.APIKeyGeolocation,
		},
		serviceStorage,
	)
	if errCrServiceGeo != nil {
		fmt.Printf(
			"creation service geo: %s\n",
			errCrServiceGeo.Error(),
		)

		os.Exit(hxerrors.OSExitForServiceIssues)
	}

	result, errCrServiceAnalytics := sanalytics.NewServiceAnalytics(
		&sanalytics.PiersNewServiceAnalytics{
			ServiceGeo: serviceGeo,
			Offsets:    &params.Offsets,
		},
		"some-site", // TODO: change
	)
	if errCrServiceAnalytics != nil {
		fmt.Printf(
			"creation service analytics: %s\n",
			errCrServiceAnalytics.Error(),
		)

		os.Exit(hxerrors.OSExitForServiceIssues)
	}

	return result,
		nil
}
