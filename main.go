package main

import (
	"agalera.eu/flight-tracking/internal/filtering/military"
	"agalera.eu/flight-tracking/internal/flight"
	"agalera.eu/flight-tracking/internal/provider/radarbox"
	"agalera.eu/flight-tracking/internal/request"
	"agalera.eu/flight-tracking/internal/zone"
	"fmt"
)

func main() {
	dataFolder := "./data"
	eu := zone.Zone{Name: "Europe", Area: zone.Bounds{Lat1: 53.827, Lng1: -10.61, Lat2: 42.163, Lng2: 63.217}}
	//provider := flightradar24.NewFlightRadarProvider(request.HttpRequestExecutor{})
	provider := radarbox.NewRadarBoxProvider(request.HttpRequestExecutor{})

	milFilter := military.NewMilitaryFlightFilter(dataFolder)
	flights := provider.GetFlights(eu)
	potentiallyMilitaryFlights := milFilter.Filter(flights)
	var flightsWithMoreDetails []flight.Flight
	for _, potentialMilitaryFlight := range potentiallyMilitaryFlights {
		moreDetailedFlight := provider.MoreDetails(potentialMilitaryFlight)
		flightsWithMoreDetails = append(flightsWithMoreDetails, moreDetailedFlight)
	}
	potentiallyMilitaryFlightsWithMoreDetails := milFilter.Filter(flightsWithMoreDetails)

	for _, potentialMilitaryFlight := range potentiallyMilitaryFlightsWithMoreDetails {
		fmt.Printf("%s\t(%s)\t[%s] -> %s\n", potentialMilitaryFlight.CallSign, potentialMilitaryFlight.AircraftType,
			potentialMilitaryFlight.Airline, provider.GetLink(potentialMilitaryFlight))
	}
}
