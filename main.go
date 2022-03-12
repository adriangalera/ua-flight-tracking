package main

import (
	"agalera.eu/flight-tracking/internal/detector"
	"agalera.eu/flight-tracking/internal/notification"
	"agalera.eu/flight-tracking/internal/provider"
	"agalera.eu/flight-tracking/internal/provider/flightradar24"
	"agalera.eu/flight-tracking/internal/provider/radarbox"
	"agalera.eu/flight-tracking/internal/request"
	"agalera.eu/flight-tracking/internal/zone"
	"time"
)

func AllProviders() []provider.DataProvider {
	executor := request.HttpRequestExecutor{}
	return []provider.DataProvider{
		flightradar24.NewFlightRadarProvider(executor),
		radarbox.NewRadarBoxProvider(executor),
	}
}

func main() {
	dataFolder := "./data"
	eu := zone.Zone{Name: "Europe centered in Kiev", Area: zone.Bounds{Lat1: 55.594, Lng1: -5.48, Lat2: 44.389, Lng2: 68.344}}
	notifiers := []notification.Notifier{&notification.ConsoleNotifier{}}
	det := detector.NewMilitaryFlightDetector(notifiers, AllProviders(), dataFolder)
	for {
		det.DetectMilitaryFlights(eu)
		time.Sleep(30 * time.Second)
	}
}
