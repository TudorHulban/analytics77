package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"

	"github.com/tudorhulban/analytics77/cmd"
	"github.com/tudorhulban/analytics77/fixtures"
	"github.com/tudorhulban/analytics77/infra/initialization"
	"github.com/tudorhulban/hxerrors"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(
			"Error: Please provide an IP address as the first argument.",
		)
		fmt.Println(
			"Usage: go run main.go <IP_ADDRESS>",
		)

		os.Exit(
			hxerrors.OSExitForApplicationIssues,
		)
	}

	ipRaw := os.Args[1]

	ip := net.ParseIP(ipRaw)
	if ip == nil {
		fmt.Printf(
			"Error: '%s' is not a valid IP address.\n",
			ipRaw,
		)

		os.Exit(
			hxerrors.OSExitForApplicationIssues,
		)
	}

	configRaw := initialization.Configuration(cmd.PathConfig)

	configSiteName, errParseSite := extractSiteName(configRaw)
	if errParseSite != nil {
		fmt.Printf(
			"error extract configuration: %s\n",
			errParseSite.Error(),
		)

		os.Exit(
			hxerrors.OSExitForConfigurationIssues,
		)
	}

	configSocket, errParse := extractServerSocketConfig(configRaw)
	if errParse != nil {
		fmt.Printf(
			"error extract configuration: %s\n",
			errParse.Error(),
		)

		os.Exit(
			hxerrors.OSExitForConfigurationIssues,
		)
	}

	connClient, errListener := net.Dial( //nolint:noctx
		"tcp",
		configSocket,
	)
	if errListener != nil {
		fmt.Printf(
			"error connecting server: %s\n",
			errListener.Error(),
		)

		os.Exit(
			hxerrors.OSExitForConfigurationIssues,
		)
	}

	request := fixtures.NewRequests(configSiteName, ip.String())

	if errTransmit := gob.NewEncoder(connClient).Encode(&request); errTransmit != nil {
		fmt.Printf(
			"error encoding request: %s\n",
			errTransmit.Error(),
		)

		os.Exit(
			hxerrors.OSExitForInfrastructureIssues,
		)
	}

	_ = connClient.Close()
}
