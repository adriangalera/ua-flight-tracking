package detector

import (
	"agalera.eu/flight-tracking/internal/filtering"
	"agalera.eu/flight-tracking/internal/filtering/military"
	"agalera.eu/flight-tracking/internal/flight"
	"agalera.eu/flight-tracking/internal/notification"
	"agalera.eu/flight-tracking/internal/provider"
	"agalera.eu/flight-tracking/internal/zone"
)

type MilitaryFlightDetector struct {
	militaryFlightsByCallSign map[string]flight.Flight
	notifiers                 []notification.Notifier
	providers                 []provider.DataProvider
	militaryFilter            filtering.Filter
}

func NewMilitaryFlightDetector(notifiers []notification.Notifier, providers []provider.DataProvider, dataFolder string) MilitaryFlightDetector {
	return MilitaryFlightDetector{
		militaryFlightsByCallSign: map[string]flight.Flight{},
		notifiers:                 notifiers,
		providers:                 providers,
		militaryFilter:            military.NewMilitaryFlightFilter(dataFolder),
	}
}

func (m *MilitaryFlightDetector) DetectMilitaryFlights(z zone.Zone) {
	var parsedMilitaryFlights []flight.Flight
	for _, flighProvider := range m.providers {
		flights := flighProvider.GetFlights(z)
		var flightsWithMoreDetails []flight.Flight
		potentiallyMilitaryFlights := m.militaryFilter.Filter(flights)
		for _, potentialMilitaryFlight := range potentiallyMilitaryFlights {
			moreDetailedFlight := flighProvider.MoreDetails(potentialMilitaryFlight)
			flightsWithMoreDetails = append(flightsWithMoreDetails, moreDetailedFlight)
			potentiallyMilitaryFlightsWithMoreDetails := m.militaryFilter.Filter(flightsWithMoreDetails)
			for _, militaryFlight := range potentiallyMilitaryFlightsWithMoreDetails {
				militaryFlight.Link = flighProvider.GetLink(militaryFlight)
				parsedMilitaryFlights = append(parsedMilitaryFlights, militaryFlight)
			}
		}
	}

	m.notifyNewFlightsInTheZone(parsedMilitaryFlights, z)
	m.notifyFlightsOutOfTheZone(parsedMilitaryFlights, z)
}

func (m *MilitaryFlightDetector) notifyNewFlightsInTheZone(currentFlights []flight.Flight, z zone.Zone) {
	for _, f := range currentFlights {
		_, isPresent := m.militaryFlightsByCallSign[f.CallSign]
		if !isPresent {
			m.militaryFlightsByCallSign[f.CallSign] = f
			for _, notifier := range m.notifiers {
				notifier.NotifyFlightIn(f, z)
			}
		}
	}
}

func (m *MilitaryFlightDetector) notifyFlightsOutOfTheZone(flights []flight.Flight, z zone.Zone) {
	currentFlights := map[string]flight.Flight{}
	for _, f := range flights {
		currentFlights[f.CallSign] = f
	}
	var flightsToDelete []flight.Flight
	for callSign, storedFlight := range m.militaryFlightsByCallSign {
		_, isPresent := currentFlights[callSign]
		if !isPresent {
			flightsToDelete = append(flightsToDelete, storedFlight)
		}
	}
	for _, flightToDelete := range flightsToDelete {
		delete(m.militaryFlightsByCallSign, flightToDelete.CallSign)
		for _, notifier := range m.notifiers {
			notifier.NotifyFlightOut(flightToDelete, z)
		}
	}
}
