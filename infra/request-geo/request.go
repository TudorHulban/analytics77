package requestgeo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/tudorhulban/analytics77/domain/analytics"
)

type ParamsGetLocationByIP struct {
	Client *http.Client

	APIKey    string
	IPAddress string
}

func GetLocationByIP(params *ParamsGetLocationByIP) (*analytics.GeoIP, error) {
	// 1. Construct the URL safely
	baseURL, errParse := url.Parse("https://api.ipgeolocation.io/v3/ipgeo")
	if errParse != nil {
		return nil,
			fmt.Errorf("failed parsing base URL: %w", errParse)
	}

	urlValues := url.Values{}
	urlValues.Add("ip", params.IPAddress)
	urlValues.Add("apiKey", params.APIKey)
	baseURL.RawQuery = urlValues.Encode()

	// 2. Execute the HTTP GET request
	resp, errGet := params.Client.Get(baseURL.String()) //nolint:noctx
	if errGet != nil {
		return nil,
			fmt.Errorf("http request failed: %w", errGet)
	}
	defer resp.Body.Close()

	// 3. Handle non-200 responses safely
	if resp.StatusCode != http.StatusOK {
		return nil,
			fmt.Errorf(
				"geoapify API returned status: %d",
				resp.StatusCode,
			)
	}

	// 4. Decode the JSON stream directly into the struct (more efficient than io.ReadAll)
	var geoData analytics.GeoIP

	if errDecoder := json.NewDecoder(resp.Body).Decode(&geoData); errDecoder != nil {
		return nil,
			fmt.Errorf("failed decoding json response: %w", errDecoder)
	}

	return &geoData, nil
}
