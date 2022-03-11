package radarbox

import (
	"agalera.eu/flight-tracking/internal/flight"
	"agalera.eu/flight-tracking/internal/request"
	"agalera.eu/flight-tracking/internal/zone"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestParseMoreDetailsResponse(t *testing.T) {

}

func TestGetFlights(t *testing.T) {
	content, err := ioutil.ReadFile("./data/radarbox.json")
	assert.Nil(t, err, "Error reading radarbox.json file with error %s", err)
	provider := NewRadarBoxProvider(request.MockRequestExecutor{Response: string(content)})
	bounds := zone.Bounds{Lat1: 1, Lat2: 2, Lng1: 3, Lng2: 4}
	testZone := zone.Zone{
		Name: "test",
		Area: bounds,
	}
	parsedFlights := provider.GetFlights(testZone)
	flight1 := flight.Flight{
		Id:       "1747906115",
		CallSign: "LH8433",
	}
	flight2 := flight.Flight{
		Id:       "1748215643",
		CallSign: "FR3619",
	}
	assert.Equal(t, flight1, parsedFlights[0], "First aircraft does not match expected")
	assert.Equal(t, flight2, parsedFlights[len(parsedFlights)-1], "Last aircraft does not match expected")
}

func TestGetLink(t *testing.T) {
	flightToLink := flight.Flight{
		Id:           "1748155607",
		CallSign:     "ETD913",
		AircraftType: "Boeing 777-FFX",
		Airline:      "Etihad Airways",
	}
	provider := NewRadarBoxProvider(request.MockRequestExecutor{})
	link := provider.GetLink(flightToLink)
	expectedLink := fmt.Sprintf("https://www.radarbox.com/flight/%s", flightToLink.CallSign)
	assert.Equal(t, expectedLink, link, "link does not match the expected")
}

func TestMoreDetails(t *testing.T) {
	content, err := ioutil.ReadFile("./data/radarbox-flight.json")
	assert.Nil(t, err, "Error reading radarbox-flight.json file with error %s", err)
	provider := NewRadarBoxProvider(request.MockRequestExecutor{Response: string(content)})
	initialFlight := flight.Flight{
		Id:       "1748155607",
		CallSign: "ETD913",
	}
	parsedFlight := provider.MoreDetails(initialFlight)
	expectedFlight := flight.Flight{
		Id:           "1748155607",
		CallSign:     "ETD913",
		AircraftType: "Boeing 777-FFX",
		Airline:      "Etihad Airways",
	}
	assert.Equal(t, expectedFlight, parsedFlight, "flight does not match the expected")
}
