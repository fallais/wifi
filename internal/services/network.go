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

type networkService struct {
}

//------------------------------------------------------------------------------
// Factory
//------------------------------------------------------------------------------

// NewNetworkService returns a new NetworkService
func NewNetworkService() NetworkService {
	return &networkService{}
}

//------------------------------------------------------------------------------
// Services
//------------------------------------------------------------------------------

// GetWiFiInfo returns the WiFi connection information.
func (service *networkService) ListAvailableNetworks() ([]*models.Network, error) {
	cmd := exec.Command("iwlist", "scan")
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(stdout)))

	// Create the network
	networks, err := models.NewNetworkFromScanner(scanner)
	if err != nil {
		log.Printf("error creating network: %v", err)
	}

	return networks, nil

}
