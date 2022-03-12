package detector

import (
	"agalera.eu/flight-tracking/internal/flight"
	"agalera.eu/flight-tracking/internal/notification"
	"agalera.eu/flight-tracking/internal/provider"
	"agalera.eu/flight-tracking/internal/zone"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testZone = zone.Zone{Name: "test", Area: zone.Bounds{}}
var flight1 = flight.Flight{CallSign: "NCHO23"}
var flight2 = flight.Flight{CallSign: "VIPER21"}
var flight3 = flight.Flight{CallSign: "DUKE23"}
var dataFolder = "../../data"

func TestDetectFlightInWhenEmpty(t *testing.T) {
	notifier := &mockNotifier{}
	notifiers := []notification.Notifier{notifier}
	providers := []provider.DataProvider{&mockProvider{nextFlights: []flight.Flight{flight1}}}
	detector := NewMilitaryFlightDetector(notifiers, providers, dataFolder)
	detector.DetectMilitaryFlights(testZone)
	assert.Equal(t, flight1, notifier.flightIn, "flight not notified")
}

func TestDetectFlightInWhenNotEmpty(t *testing.T) {
	notifier := &mockNotifier{}
	notifiers := []notification.Notifier{notifier}
	mockProvider := mockProvider{nextFlights: []flight.Flight{flight1}}
	providers := []provider.DataProvider{&mockProvider}
	detector := NewMilitaryFlightDetector(notifiers, providers, dataFolder)
	detector.DetectMilitaryFlights(testZone)
	mockProvider.nextFlights = []flight.Flight{flight2}
	detector.DetectMilitaryFlights(testZone)
	assert.Equal(t, flight2, notifier.flightIn, "flight not notified")
}

func TestDetectFlightOut(t *testing.T) {
	notifier := &mockNotifier{}
	notifiers := []notification.Notifier{notifier}
	mockProvider := mockProvider{nextFlights: []flight.Flight{flight1, flight2, flight3}}
	providers := []provider.DataProvider{&mockProvider}
	detector := NewMilitaryFlightDetector(notifiers, providers, dataFolder)

	detector.DetectMilitaryFlights(testZone)
	mockProvider.nextFlights = []flight.Flight{flight1, flight3}
	detector.DetectMilitaryFlights(testZone)

	assert.Equal(t, flight2, notifier.flightOut, "flight not notified")
}

type mockNotifier struct {
	flightIn  flight.Flight
	flightOut flight.Flight
}

func (m *mockNotifier) NotifyFlightIn(flight flight.Flight, zone zone.Zone) {
	m.flightIn = flight
}

func (m *mockNotifier) NotifyFlightOut(flight flight.Flight, zone zone.Zone) {
	m.flightOut = flight
}

type mockProvider struct {
	nextFlights []flight.Flight
}

func (m *mockProvider) GetFlights(zone zone.Zone) []flight.Flight {
	return m.nextFlights
}

func (m *mockProvider) GetLink(flight flight.Flight) string {
	return ""
}

func (m *mockProvider) MoreDetails(flight flight.Flight) flight.Flight {
	return flight
}
