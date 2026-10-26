package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/tudorhulban/analytics77/cmd"
	"github.com/tudorhulban/analytics77/domain/site"
)

type configuration struct {
	portRPC  string
	portHTTP string

	nameLogfile    string
	utcOffsetHours int64
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

	initialization, exists := raw["initialization"].(map[string]any)
	if !exists {
		return nil,
			errors.New(
				"invalid or missing initialization section",
			)
	}

	rawSites, exists := initialization["sites"].([]any)
	if !exists {
		return nil, errors.New("invalid or missing sites configuration")
	}

	var sites []*site.Site

	for _, rawSite := range rawSites {
		siteMap, ok := rawSite.(map[string]any)
		if !ok {
			return nil, errors.New("invalid site entry format")
		}

		item, errCr := site.NewSite(siteMap)
		if errCr != nil {
			return nil, errCr
		}
		sites = append(sites, item)
	}

	offset, couldCastOffset := initialization["utcOffsetHours"].(float64)
	if !couldCastOffset {
		return nil,
			fmt.Errorf(
				"invalid initialization value as %v",
				initialization["utcOffsetHours"],
			)
	}

	return &configuration{
			portRPC:     strconv.Itoa(int(portRPC)),
			portHTTP:    strconv.Itoa(int(portHTTP)),
			nameLogfile: nameLogfile,

			utcOffsetHours: int64(offset),
		},
		nil
}
