package services

import "wifi/internal/models"

type WifiConnectionService interface {
	GetWiFiInfo() (*models.WifiConnection, error)
}

type NetworkService interface {
	ListAvailableNetworks() ([]*models.Network, error)
}
