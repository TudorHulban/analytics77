package initialization

import (
	"log"
	"os"

	"github.com/tudorhulban/analytics77/helpers"
	"github.com/tudorhulban/hxerrors"
)

func Configuration(path string) map[string]any {
	result, errParse := helpers.ParseJSONFile(path)
	if errParse != nil {
		log.Printf(
			"configuration error: %s",
			errParse.Error(),
		)

		os.Exit(hxerrors.OSExitForConfigurationIssues)
	}

	return result
}
