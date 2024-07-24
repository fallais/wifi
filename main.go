package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"
	"wifi/internal/models"

	"github.com/rivo/tview"
)

const Refresh = 1 * time.Second

// Function to get WiFi stats
func getWiFiInfo() (*models.WifiConnection, error) {
	cmd := exec.Command("iwconfig")
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(stdout)))

	// Create the network
	wifiConnection, err := models.NewFromScanner(scanner)
	if err != nil {
		log.Printf("error creating network: %v", err)
	}

	return wifiConnection, nil
}

// Function to scan WiFi networks
func scanWiFiNetworks() ([]map[string]string, error) {
	cmd := exec.Command("iwlist", "scan")
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	networks := []map[string]string{}
	scanner := bufio.NewScanner(strings.NewReader(string(stdout)))

	reSSID := regexp.MustCompile(`ESSID:"(.*?)"`)
	reSignal := regexp.MustCompile(`Signal level=(-?\d+) dBm`)
	reQuality := regexp.MustCompile(`Quality=(\d+)/(\d+)`)

	var currentNetwork map[string]string

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Cell") {
			if currentNetwork != nil {
				networks = append(networks, currentNetwork)
			}
			currentNetwork = make(map[string]string)
		}
		if match := reSSID.FindStringSubmatch(line); match != nil {
			currentNetwork["ESSID"] = match[1]
		}
		if match := reSignal.FindStringSubmatch(line); match != nil {
			currentNetwork["Signal Level"] = match[1]
		}
		if match := reQuality.FindStringSubmatch(line); match != nil {
			currentNetwork["Link Quality"] = fmt.Sprintf("%s/%s", match[1], match[2])
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

func main() {
	app := tview.NewApplication()

	// Create text views for WiFi stats and networks
	wifiStatsView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	wifiNetworksView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	// Create a pages view to switch between tabs
	pages := tview.NewPages().
		AddPage("WiFi Stats", wifiStatsView, true, true).
		AddPage("WiFi Networks", wifiNetworksView, true, false)

	// Create a list to serve as the tab selector
	tabList := tview.NewList().
		AddItem("WiFi Stats", "View WiFi stats", '1', func() {
			pages.SwitchToPage("WiFi Stats")
		}).
		AddItem("WiFi Networks", "View available WiFi networks", '2', func() {
			pages.SwitchToPage("WiFi Networks")
		}).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})

	// Create a flex layout to hold the tab list and the pages
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tabList, 0, 1, true).
		AddItem(pages, 0, 3, false)

	// Update WiFi stats periodically
	go func() {
		for {
			wifiConnection, err := getWiFiInfo()
			if err != nil {
				fmt.Fprintf(wifiStatsView, "[red]Error: %s", err)
				app.Draw()
				return
			}

			wifiStatsView.Clear()
			fmt.Fprintf(wifiStatsView, "[yellow]WiFi Information\n")
			fmt.Fprintf(wifiStatsView, "[yellow]ESSID: [white]%s\n", wifiConnection.ESSID)
			fmt.Fprintf(wifiStatsView, "[yellow]Signal Level: [white]%s dBm\n", wifiConnection.SignalLevel)
			fmt.Fprintf(wifiStatsView, "[yellow]Link Quality: [white]%s\n", wifiConnection.LinkQuality)
			fmt.Fprintf(wifiStatsView, "[yellow]Tx Retries: [white]%s\n", wifiConnection.TxRetries)

			time.Sleep(Refresh)
		}
	}()

	// Update WiFi networks periodically
	go func() {
		for {
			networks, err := scanWiFiNetworks()
			if err != nil {
				fmt.Fprintf(wifiNetworksView, "[red]Error: %s", err)
				app.Draw()
				return
			}

			wifiNetworksView.Clear()
			fmt.Fprintf(wifiNetworksView, "[yellow]Available WiFi Networks\n\n")
			for _, network := range networks {
				fmt.Fprintf(wifiNetworksView, "[yellow]ESSID: [white]%s\n", network["ESSID"])
				fmt.Fprintf(wifiNetworksView, "[yellow]Signal Level: [white]%s dBm\n", network["Signal Level"])
				fmt.Fprintf(wifiNetworksView, "[yellow]Link Quality: [white]%s\n\n", network["Link Quality"])
			}

			time.Sleep(Refresh)
		}
	}()

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
