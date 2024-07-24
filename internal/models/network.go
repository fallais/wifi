package models

import (
	"bufio"
	"fmt"
	"regexp"
)

var reSSID = regexp.MustCompile(`ESSID:"(.*?)"`)
var reSignal = regexp.MustCompile(`Signal level=(-?\d+) dBm`)
var reQuality = regexp.MustCompile(`Link Quality=(\d+)/(\d+)`)
var reRetries = regexp.MustCompile(`Tx excessive retries:(\d+)`)

//------------------------------------------------------------------------------
// Structure
//------------------------------------------------------------------------------

// WifiConnection represents a Wifi connection.
type WifiConnection struct {
	ESSID       string
	SignalLevel string
	LinkQuality string
	TxRetries   string
}

//------------------------------------------------------------------------------
// Factory
//------------------------------------------------------------------------------

// NewFromScanner creates a new Network instance from a budio scanner.
func NewFromScanner(scanner *bufio.Scanner) (*WifiConnection, error) {
	wifiConnection := &WifiConnection{}

	for scanner.Scan() {
		line := scanner.Text()
		if match := reSSID.FindStringSubmatch(line); match != nil {
			wifiConnection.ESSID = match[1]
		}
		if match := reSignal.FindStringSubmatch(line); match != nil {
			wifiConnection.SignalLevel = match[1]
		}
		if match := reQuality.FindStringSubmatch(line); match != nil {
			wifiConnection.LinkQuality = fmt.Sprintf("%s/%s", match[1], match[2])
		}
		if match := reRetries.FindStringSubmatch(line); match != nil {
			wifiConnection.TxRetries = match[1]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return wifiConnection, nil
}
