package provider

import (
	"agalera.eu/flight-tracking/internal/flight"
	"agalera.eu/flight-tracking/internal/zone"
)

type DataProvider interface {
	GetFlights(zone zone.Zone) []flight.Flight

	GetLink(flight flight.Flight) string

	MoreDetails(flight flight.Flight) flight.Flight
}
