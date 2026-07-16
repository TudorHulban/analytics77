package main

import (
	"errors"
	"fmt"

	"github.com/tudorhulban/analytics77/cmd"
)

func extractConfiguration(raw map[string]any) (string, error) {
	debug, exists := raw["debug"].(map[string]any)
	if !exists {
		return "",
			errors.New(
				"invalid or missing debug configuration",
			)
	}

	// Go defaults JSON numbers to float64.
	// Type assertions return the zero-value (empty string / 0) if they fail.
	host, couldCastHost := debug["host"].(string)
	if !couldCastHost {
		return "",
			fmt.Errorf(
				"invalid host as %v",
				debug["host"],
			)
	}

	port, couldCastPort := debug[cmd.PortRPC].(float64)
	if !couldCastPort {
		return "",
			fmt.Errorf(
				"invalid port as %v",
				debug["port"],
			)
	}

	return fmt.Sprintf(
			"%s:%d",

			host,
			int(port),
		),
		nil
}
