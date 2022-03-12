package main

import (
	"agalera.eu/flight-tracking/internal/filtering/military"
	"agalera.eu/flight-tracking/internal/flight"
	"agalera.eu/flight-tracking/internal/provider"
	"agalera.eu/flight-tracking/internal/provider/flightradar24"
	"agalera.eu/flight-tracking/internal/provider/radarbox"
	"agalera.eu/flight-tracking/internal/request"
	"agalera.eu/flight-tracking/internal/zone"
	"fmt"
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
	milFilter := military.NewMilitaryFlightFilter(dataFolder)

	eu := zone.Zone{Name: "Europe", Area: zone.Bounds{Lat1: 53.827, Lng1: -10.61, Lat2: 42.163, Lng2: 63.217}}

	currentMilitaryFlights := map[string]flight.Flight{}
	for _, flighProvider := range AllProviders() {
		flights := flighProvider.GetFlights(eu)
		var flightsWithMoreDetails []flight.Flight
		potentiallyMilitaryFlights := milFilter.Filter(flights)
		for _, potentialMilitaryFlight := range potentiallyMilitaryFlights {
			moreDetailedFlight := flighProvider.MoreDetails(potentialMilitaryFlight)
			flightsWithMoreDetails = append(flightsWithMoreDetails, moreDetailedFlight)
			potentiallyMilitaryFlightsWithMoreDetails := milFilter.Filter(flightsWithMoreDetails)
			for _, militaryFlight := range potentiallyMilitaryFlightsWithMoreDetails {
				militaryFlight.Link = flighProvider.GetLink(militaryFlight)
				_, isPresent := currentMilitaryFlights[militaryFlight.CallSign]
				if !isPresent {
					currentMilitaryFlights[militaryFlight.CallSign] = militaryFlight
				}
			}
		}
	}

	for _, militaryFlight := range currentMilitaryFlights {
		fmt.Printf("%s\t%s (%s) - %s\n", militaryFlight.CallSign, militaryFlight.AircraftType, militaryFlight.Airline, militaryFlight.Link)
	}
}
