package services

import (
	"bufio"
	"log"
	"os/exec"
	"strings"

	"wifi/internal/models"
)

//------------------------------------------------------------------------------
// Structure
//------------------------------------------------------------------------------

type wifiConnectionService struct {
}

//------------------------------------------------------------------------------
// Factory
//------------------------------------------------------------------------------

// NewWifiConnectionService returns a new WifiConnectionService
func NewWifiConnectionService() WifiConnectionService {
	return &wifiConnectionService{}
}

//------------------------------------------------------------------------------
// Services
//------------------------------------------------------------------------------

// GetWiFiInfo returns the WiFi connection information.
func (service *wifiConnectionService) GetWiFiInfo() (*models.WifiConnection, error) {
	cmd := exec.Command("iwconfig")
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(stdout)))

	// Create the network
	wifiConnection, err := models.NewWifiConnectionFromScanner(scanner)
	if err != nil {
		log.Printf("error creating network: %v", err)
	}

	return wifiConnection, nil
}
