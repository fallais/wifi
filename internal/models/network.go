package models

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

var networkRegexSSID = regexp.MustCompile(`ESSID:"(.*?)"`)
var networkRegexSignal = regexp.MustCompile(`Signal level=(-?\d+) dBm`)
var networkRegexQuality = regexp.MustCompile(`Quality=(\d+)/(\d+)`)
var networkRegexMAC = regexp.MustCompile(`Address: ([\da-fA-F:]+)`)

// -----------------------------------------------------------------------------
// Structure
// -----------------------------------------------------------------------------

type Network struct {
	ESSID       string
	SignalLevel string
	MACAddress  string
	LinkQuality string
}

// -----------------------------------------------------------------------------
// Factory
// -----------------------------------------------------------------------------

// NewNetworkFromScanner creates a new Network instance from a budio scanner.
func NewNetworkFromScanner(scanner *bufio.Scanner) ([]*Network, error) {
	var networks []*Network

	var currentNetwork *Network
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Cell") {
			if currentNetwork != nil {
				networks = append(networks, currentNetwork)
			}
			currentNetwork = &Network{}
		}
		if match := networkRegexSSID.FindStringSubmatch(line); match != nil {
			currentNetwork.ESSID = match[1]
		}
		if match := networkRegexSignal.FindStringSubmatch(line); match != nil {
			currentNetwork.SignalLevel = match[1]
		}
		if match := networkRegexQuality.FindStringSubmatch(line); match != nil {
			currentNetwork.LinkQuality = fmt.Sprintf("%s/%s", match[1], match[2])
		}
		if match := networkRegexMAC.FindStringSubmatch(line); match != nil {
			currentNetwork.MACAddress = match[1]
		}
	}

	if currentNetwork != nil {
		networks = append(networks, currentNetwork)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return networks, nil
}
