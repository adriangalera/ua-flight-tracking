package flightradar24

import (
	"agalera.eu/flight-tracking/internal/flight"
	"agalera.eu/flight-tracking/internal/provider"
	"agalera.eu/flight-tracking/internal/request"
	"agalera.eu/flight-tracking/internal/zone"
	"fmt"
	"github.com/tidwall/gjson"
	"log"
)

type FlightRadarProvider struct {
	requestExecutor request.RequestExecutor
}

func NewFlightRadarProvider(requestExecutor request.RequestExecutor) provider.DataProvider {
	return &FlightRadarProvider{requestExecutor: requestExecutor}
}

func getUrl(bounds zone.Bounds) string {
	url := "https://data-cloud.flightradar24.com/zones/fcgi/feed.js"
	params := "?faa=1&satellite=1&mlat=1&flarm=1&adsb=1&gnd=0&air=1&vehicles=1&estimated=1&maxage=14400&gliders=0&stats=0"
	boundsParams := fmt.Sprintf("&bounds=%f%%2C%f%%2C%f%%2C%f", bounds.Lat1, bounds.Lat2, bounds.Lng1, bounds.Lng2)
	params = params + boundsParams
	return url + params
}

func getMoreDetailsUrl(flight flight.Flight) string {
	return fmt.Sprintf("https://data-live.flightradar24.com/clickhandler/?version=1.5&flight=%s", flight.Id)
}

func parseFlightsResponse(body string) []flight.Flight {
	var flights []flight.Flight
	keys := gjson.Get(body, "@keys")
	for _, key := range keys.Array() {
		keyStr := key.String()
		if keyStr != "full_count" && keyStr != "version" {
			flightDetailsArray := gjson.Get(body, keyStr)
			parsedFlight := flight.Flight{
				Id:       keyStr,
				CallSign: flightDetailsArray.Array()[16].String(),
			}
			flights = append(flights, parsedFlight)
		}
	}
	return flights
}

func parseMoreDetailsResponse(body string) flight.Flight {
	id := gjson.Get(body, "identification.id")
	callsign := gjson.Get(body, "identification.callsign")
	aircraftType := gjson.Get(body, "aircraft.model.text")
	airline := gjson.Get(body, "airline.name")

	return flight.Flight{
		Id:           id.String(),
		CallSign:     callsign.String(),
		AircraftType: aircraftType.String(),
		Airline:      airline.String(),
	}
}

func (f *FlightRadarProvider) GetFlights(zone zone.Zone) []flight.Flight {
	log.Printf("Query flightradar for zone %s", zone.Name)
	flightsUrl := getUrl(zone.Area)
	body := f.requestExecutor.Get(flightsUrl, map[string]string{})
	return parseFlightsResponse(body)
}

func (f *FlightRadarProvider) GetLink(flight flight.Flight) string {
	return fmt.Sprintf("https://www.flightradar24.com/%s/%s", flight.CallSign, flight.Id)
}

func (f *FlightRadarProvider) MoreDetails(flight flight.Flight) flight.Flight {
	body := f.requestExecutor.Get(getMoreDetailsUrl(flight), map[string]string{})
	parsedFlight := parseMoreDetailsResponse(body)
	parsedFlight.CallSign = flight.CallSign
	return parsedFlight
}
