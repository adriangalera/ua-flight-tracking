package config

import (
	"agalera.eu/flight-tracking/internal/zone"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNoZone(t *testing.T) {
	conf, err := ReadConfiguration("./examples/nozone.yml")
	assert.Nil(t, err, "Error reading configuration")
	expectedConf := &Configuration{Zones: nil}
	assert.Equal(t, expectedConf, conf, "Read configuration has not expected values")
}

func TestOneZone(t *testing.T) {
	conf, err := ReadConfiguration("./examples/onezone.yml")
	assert.Nil(t, err, "Error reading configuration")

	var zones []zone.Zone
	zones = append(zones, zone.Zone{
		Name: "test1",
		Area: zone.Bounds{Lat1: 1, Lng1: 2, Lat2: 1, Lng2: 2},
	})
	expectedConf := &Configuration{Zones: zones}
	assert.Equal(t, expectedConf, conf, "Read configuration has not expected values")
}

func TestTwoZone(t *testing.T) {
	conf, err := ReadConfiguration("./examples/twozones.yml")
	assert.Nil(t, err, "Error reading configuration")

	var zones []zone.Zone
	zone1 := zone.Zone{
		Name: "test1",
		Area: zone.Bounds{Lat1: 1, Lng1: 2, Lat2: 1, Lng2: 2},
	}
	zone2 := zone.Zone{
		Name: "test2",
		Area: zone.Bounds{Lat1: 3, Lng1: 4, Lat2: 3, Lng2: 4},
	}
	zones = append(zones, zone1)
	zones = append(zones, zone2)
	expectedConf := &Configuration{Zones: zones}
	assert.Equal(t, expectedConf, conf, "Read configuration has not expected values")
}

func TestErrorReading(t *testing.T) {
	_, err := ReadConfiguration("/tmp/nonexisting")
	assert.NotNil(t, err, "Expected error did not happen")
}
