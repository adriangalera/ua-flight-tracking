package flightradar24

import (
	"agalera.eu/flight-tracking/internal/flight"
	"agalera.eu/flight-tracking/internal/request"
	"agalera.eu/flight-tracking/internal/zone"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestGetUrl(t *testing.T) {
	bounds := zone.Bounds{Lat1: 1, Lat2: 2, Lng1: 3, Lng2: 4}
	url := getUrl(bounds)
	expectedUrl := "https://data-cloud.flightradar24.com/zones/fcgi/feed.js?faa=1&satellite=1&mlat=1&flarm=1&adsb=1&gnd=0&air=1&vehicles=1&estimated=1&maxage=14400&gliders=0&stats=0&bounds=1.000000%2C2.000000%2C3.000000%2C4.000000"
	assert.Equal(t, expectedUrl, url, "URL not matching the expected")
}

func TestGetMoreDetails(t *testing.T) {
	content, err := ioutil.ReadFile("./data/flightDetails.json")
	assert.Nil(t, err, "Error reading flightDetails.json file with error %s", err)
	provider := NewFlightRadarProvider(request.MockRequestExecutor{Response: string(content)})
	testFlight := flight.Flight{
		Id:       "2b19efdb",
		CallSign: "LAGR241",
	}
	expectedFlight := flight.Flight{
		Id:           "2b19efdb",
		CallSign:     "LAGR241",
		AircraftType: "Boeing KC-135R Stratotanker",
		Airline:      "United States - US Air Force (USAF)",
	}
	parsedFlight := provider.MoreDetails(testFlight)
	assert.Equal(t, expectedFlight, parsedFlight, "Parsed flight does not match the expected")
}

func TestGetMoreDetailsUrl(t *testing.T) {
	url := getMoreDetailsUrl(flight.Flight{Id: "1"})
	expectedUrl := "https://data-live.flightradar24.com/clickhandler/?version=1.5&flight=1"
	assert.Equal(t, expectedUrl, url, "more details URL do not match")
}

func TestGetLink(t *testing.T) {
	provider := NewFlightRadarProvider(request.MockRequestExecutor{})
	flightId := "1"
	flightCallSign := "a"
	link := provider.GetLink(flight.Flight{Id: flightId, CallSign: flightCallSign})
	expectedLink := fmt.Sprintf("https://www.flightradar24.com/%s/%s", flightCallSign, flightId)
	assert.Equal(t, expectedLink, link, "link do not match")
}

func TestGetFlights(t *testing.T) {
	content, err := ioutil.ReadFile("./data/fr24.json")
	assert.Nil(t, err, "Error reading fr24.json file with error %s", err)
	provider := NewFlightRadarProvider(request.MockRequestExecutor{Response: string(content)})
	bounds := zone.Bounds{Lat1: 1, Lat2: 2, Lng1: 3, Lng2: 4}
	testZone := zone.Zone{
		Name: "test",
		Area: bounds,
	}
	flights := provider.GetFlights(testZone)
	firstFlight := flight.Flight{
		Id:       "2b166e3c",
		CallSign: "C25B",
	}
	lastFlight := flight.Flight{
		Id:       "2b16af97",
		CallSign: "PHWMA",
	}
	assert.Equal(t, firstFlight, flights[0], "First flight does not match expected")
	assert.Equal(t, lastFlight, flights[len(flights)-1], "Last flight does not match expected")
}


