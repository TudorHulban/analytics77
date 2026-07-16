package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/tudorhulban/analytics77/cmd"
)

type configuration struct {
	portRPC  string
	portHTTP string

	nameLogfile string
}

func extractConfiguration(raw map[string]any) (*configuration, error) {
	server, exists := raw["server"].(map[string]any)
	if !exists {
		return nil,
			errors.New(
				"invalid or missing server configuration",
			)
	}

	// Go defaults JSON numbers to float64.
	// Type assertions return the zero-value (empty string / 0) if they fail.
	portRPC, couldCastPortRPC := server[cmd.PortRPC].(float64)
	if !couldCastPortRPC {
		return nil,
			fmt.Errorf(
				"invalid port RPC as %v",
				server[cmd.PortRPC],
			)
	}

	portHTTP, couldCastPortHTTP := server[cmd.PortHTTP].(float64)
	if !couldCastPortHTTP {
		return nil,
			fmt.Errorf(
				"invalid port HTTP as %v",
				server[cmd.PortHTTP],
			)
	}

	nameLogfile, couldCastNameLogfile := server["logfile"].(string)
	if !couldCastNameLogfile {
		return nil,
			fmt.Errorf(
				"invalid name log file as %v",
				server["logfile"],
			)
	}

	return &configuration{
			portRPC:     strconv.Itoa(int(portRPC)),
			portHTTP:    strconv.Itoa(int(portHTTP)),
			nameLogfile: nameLogfile,
		},
		nil
}
