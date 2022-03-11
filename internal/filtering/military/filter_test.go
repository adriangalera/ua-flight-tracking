package military

import (
	"agalera.eu/flight-tracking/internal/flight"
	"github.com/stretchr/testify/assert"
	"testing"
)

var militaryFlight = flight.Flight{
	CallSign:     "NCHO125",
	AircraftType: "E8",
}
var nonMilitaryFlight = flight.Flight{
	CallSign:     "FIN122",
	AircraftType: "A359",
}
var confusingFlight = flight.Flight{
	CallSign:     "TRA123",
	AircraftType: "A359",
}

func TestIsMilitary(t *testing.T) {
	filter := MilitaryFlightFilter{knownMilitaryCodes: []string{"NCHO"}}
	assert.True(t, filter.isMilitary(militaryFlight), "Military flight not detected")
	assert.False(t, filter.isMilitary(nonMilitaryFlight), "Non military flight detected as military")
}

func TestFindMilitaryFlights(t *testing.T) {
	filter := MilitaryFlightFilter{knownMilitaryCodes: []string{"NCHO"}}
	allFlights := []flight.Flight{militaryFlight, nonMilitaryFlight}
	militaryFlights := filter.Filter(allFlights)
	assert.Equal(t, 1, len(militaryFlights), "More than one flight detected or zero")
	assert.Equal(t, militaryFlight, militaryFlights[0], "The detected flight does not match the expected")
}

func TestReadKnownCodes(t *testing.T) {
	dataFolder := "../../../data"
	knownMilitaryCodes := getKnownMilitaryCodes(dataFolder)
	knownCivilCodes := getKnownCivilCodes(dataFolder)
	assert.NotNil(t, knownMilitaryCodes, "Known codes is nil")
	assert.True(t, len(knownMilitaryCodes) > 0, "0 known codes")
	assert.NotNil(t, knownCivilCodes, "Known codes is nil")
	assert.True(t, len(knownCivilCodes) > 0, "0 known codes")

	filter := NewMilitaryFlightFilter(dataFolder)
	assert.NotNil(t, filter, "filter was not created properly")
}

func TestCivilFlightNotDetected(t *testing.T) {
	filter := MilitaryFlightFilter{knownMilitaryCodes: []string{"TR"}, knownCivilCodes: []string{"TRA"}}
	allFlights := []flight.Flight{confusingFlight}
	militaryFlights := filter.Filter(allFlights)
	assert.Equal(t, 0, len(militaryFlights), "Detected civil flight as military")
}
