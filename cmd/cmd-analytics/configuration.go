package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/tudorhulban/analytics77/cmd"
	"github.com/tudorhulban/analytics77/domain/site"
)

type configuration struct {
	portRPC  string
	portHTTP string

	nameLogfile string

	sites []*site.Site
}

// String implements the fmt.Stringer interface for configuration.
func (c *configuration) String() string {
	if c == nil {
		return "<nil configuration>"
	}

	var sb strings.Builder

	sb.WriteString("Configuration {\n")
	fmt.Fprintf(&sb, "  portRPC: %s\n", c.portRPC)
	fmt.Fprintf(&sb, "  portHTTP: %s\n", c.portHTTP)
	fmt.Fprintf(&sb, "  nameLogfile: %s\n", c.nameLogfile)

	sb.WriteString("  sites: [\n")

	for ix, sit := range c.sites {
		if sit == nil {
			sb.WriteString("    <nil site>,\n")

			continue
		}

		fmt.Fprintf(&sb, "    [%d] Site {\n", ix)
		fmt.Fprintf(&sb, "      Name: %v\n", sit.Name)
		fmt.Fprintf(&sb, "      OffsetUTCHours: %d\n", sit.OffsetUTCHours)
		fmt.Fprintf(&sb, "      TimestampDSTWinter: %d\n", sit.TimestampDSTWinter)
		fmt.Fprintf(&sb, "      TimestampDSTSpring: %d\n", sit.TimestampDSTSpring)

		sb.WriteString("    },\n")
	}

	sb.WriteString("  ]\n")
	sb.WriteString("}")

	return sb.String()
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

	return &configuration{
			portRPC:     strconv.Itoa(int(portRPC)),
			portHTTP:    strconv.Itoa(int(portHTTP)),
			nameLogfile: nameLogfile,

			sites: sites,
		},
		nil
}
