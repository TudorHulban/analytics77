package main

import (
	"errors"
	"fmt"

	"github.com/tudorhulban/analytics77/cmd"
)

func extractServerSocketConfig(raw map[string]any) (string, error) {
	debug, exists := raw["debug"].(map[string]any)
	if !exists {
		return "",
			errors.New(
				"invalid or missing debug configuration",
			)
	}

	// Go defaults JSON numbers to float64.
	// Type assertions return the zero-value (empty string / 0) if they fail.
	host, couldCastHost := debug[_HostServer].(string)
	if !couldCastHost {
		return "",
			fmt.Errorf(
				"invalid host as %v",
				debug[_HostServer],
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

func extractSiteName(raw map[string]any) (string, error) {
	initialization, exists := raw[_NameSectionInit].(map[string]any)
	if !exists {
		return "",
			fmt.Errorf(
				"invalid or missing %s configuration",
				_NameSectionInit,
			)
	}

	nameLocal, couldCastSiteName := initialization[_NameSite].(string)
	if !couldCastSiteName {
		return "",
			fmt.Errorf(
				"invalid site / local name as %v",
				initialization[_NameSite],
			)
	}

	return nameLocal, nil
}
