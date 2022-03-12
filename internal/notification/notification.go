package notification

import (
	"agalera.eu/flight-tracking/internal/flight"
	"agalera.eu/flight-tracking/internal/zone"
	"fmt"
)

type Notifier interface {
	NotifyFlightIn(flight flight.Flight, zone zone.Zone)
	NotifyFlightOut(flight flight.Flight, zone zone.Zone)
}

type ConsoleNotifier struct {
}

func (c *ConsoleNotifier) NotifyFlightIn(flight flight.Flight, zone zone.Zone) {
	fmt.Printf("IN [%s] %s %s %s %s\n", zone.Name, flight.CallSign, flight.AircraftType, flight.Airline, flight.Link)
}

func (c *ConsoleNotifier) NotifyFlightOut(flight flight.Flight, zone zone.Zone) {
	fmt.Printf("OUT [%s] %s %s %s %s\n", zone.Name, flight.CallSign, flight.AircraftType, flight.Airline, flight.Link)
}
