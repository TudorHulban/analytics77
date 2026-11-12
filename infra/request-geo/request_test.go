package requestgeo

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/tudorhulban/analytics77/cmd"
)

func TestGetLocationByIP(t *testing.T) {
	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}

	apiKey := os.Getenv(cmd.OSAPIGeolocation)
	targetIP := "82.76.117.202"

	location, errGeo := GetLocationByIP(
		&ParamsGetLocationByIP{
			Client:    &httpClient,
			APIKey:    apiKey,
			IPAddress: targetIP,
		},
	)
	if errGeo != nil {
		fmt.Printf("Error: %v\n", errGeo)

		return
	}

	fmt.Println(*location)
}
