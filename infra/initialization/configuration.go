package initialization

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/tudorhulban/analytics77/cmd"
	"github.com/tudorhulban/analytics77/domain/site"
	"github.com/tudorhulban/analytics77/helpers"
	"github.com/tudorhulban/hxerrors"
)

type Configuration struct {
	PortRPC  string
	PortHTTP string

	NameLogfile string

	Sites []*site.Site
}

func NewConfiguration(rawConfig map[string]any) (*Configuration, error) {
	server, exists := rawConfig["server"].(map[string]any)
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

	initialization, exists := rawConfig["initialization"].(map[string]any)
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

		// fmt.Println(siteMap)

		item, errCr := site.NewSite(siteMap)
		if errCr != nil {
			return nil, errCr
		}

		sites = append(sites, item)
	}

	return &Configuration{
			PortRPC:     strconv.Itoa(int(portRPC)),
			PortHTTP:    strconv.Itoa(int(portHTTP)),
			NameLogfile: nameLogfile,

			Sites: sites,
		},
		nil
}

// String implements the fmt.Stringer interface for configuration.
func (c *Configuration) String() string {
	if c == nil {
		return "<nil configuration>"
	}

	var sb strings.Builder
	sb.Grow(512) // Pre-allocate for better performance

	sb.WriteString("Configuration {\n")
	fmt.Fprintf(&sb, "  portRPC: %s\n", c.PortRPC)
	fmt.Fprintf(&sb, "  portHTTP: %s\n", c.PortHTTP)
	fmt.Fprintf(&sb, "  nameLogfile: %s\n", c.NameLogfile)

	sb.WriteString("  sites: [\n")

	for ix, sit := range c.Sites {
		if sit == nil {
			sb.WriteString("    <nil site>,\n")
			continue
		}

		fmt.Fprintf(&sb, "    [%d] Site {\n", ix)
		fmt.Fprintf(&sb, "      name: %s\n", sit.Name)

		// Indent the TimestampOffsets stringer output using SplitSeq (zero allocations for the slice)
		offsetsStr := sit.TimestampOffsets.String()
		for line := range strings.SplitSeq(strings.TrimSpace(offsetsStr), "\n") {
			fmt.Fprintf(&sb, "      %s\n", line)
		}

		sb.WriteString("    },\n")
	}

	sb.WriteString("  ]\n")
	sb.WriteString("}")

	return sb.String()
}

func GetConfigurationRaw(path string) map[string]any {
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
