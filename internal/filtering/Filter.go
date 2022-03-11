package filtering

import "agalera.eu/flight-tracking/internal/flight"

type Filter interface {
	Filter(flights []flight.Flight) []flight.Flight
}
