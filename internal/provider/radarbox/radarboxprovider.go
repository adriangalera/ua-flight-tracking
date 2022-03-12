package radarbox

import (
	"agalera.eu/flight-tracking/internal/flight"
	"agalera.eu/flight-tracking/internal/provider"
	"agalera.eu/flight-tracking/internal/request"
	"agalera.eu/flight-tracking/internal/zone"
	"fmt"
	"github.com/tidwall/gjson"
	"log"
	"sort"
)

type RadarboxProvider struct {
	requestExecutor request.RequestExecutor
}

func NewRadarBoxProvider(executor request.RequestExecutor) provider.DataProvider {
	return RadarboxProvider{requestExecutor: executor}
}

func headers() map[string]string {
	return map[string]string{
		"authority":       "data.rb24.com",
		"accept":          "application/json, text/plain, */*",
		"user-agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36",
		"origin":          "https://www.radarbox.com",
		"referer":         "https://www.radarbox.com",
		"accept-language": "es,ca;q=0.9,en-GB;q=0.8,en-US;q=0.7,en;q=0.6",
	}
}

func getFlightsUrl(bounds zone.Bounds) string {
	return fmt.Sprintf("https://data.rb24.com/live?aircraft=&airport=&fn=&far=&fms=&zoom=6&flightid="+
		"&bounds=%f,%f,%f,%f"+
		"&timestamp=1647091930737&designator=iata&showLastTrails=true&ff=false&os=web&adsb=true"+
		"&adsbsat=true&asdi=true&ocea=true&mlat=true&sate=true&uat=true&hfdl=true&esti=true&asdex=true&flarm=true"+
		"&aust=true&diverted=false&delayed=false&isga=false&ground=true&onair=true&blocked=false"+
		"&station=&class[]=%%3F&class[]=A&class[]=B&class[]=C&class[]=G&class[]=H&class[]=M&airline=&route=&country=",
		bounds.Lat1, bounds.Lng2, bounds.Lat2, bounds.Lng1,
	)
}

func parseFlightsResponse(body string) []flight.Flight {
	rootElement := gjson.Get(body, "@this")
	if !rootElement.Array()[0].Exists() {
		log.Printf("Cannot parse radarboxprovider!. Received body: %s", body)
		return []flight.Flight{}
	}
	flights := rootElement.Array()[0].Map()
	flightIds := make([]string, 0, len(flights))
	for flightId := range flights {
		flightIds = append(flightIds, flightId)
	}
	sort.Strings(flightIds)

	var parsedFlights []flight.Flight
	for _, flightId := range flightIds {
		curFlightData := flights[flightId]
		parsedFlights = append(parsedFlights, flight.Flight{
			Id:       flightId,
			CallSign: curFlightData.Array()[0].String(),
		})
	}
	return parsedFlights
}

func getMoreDetailsUrl(flight flight.Flight) string {
	return fmt.Sprintf("https://data.rb24.com/live-flight-info?fid=%s&locale=en", flight.Id)
}

func (r RadarboxProvider) GetFlights(zone zone.Zone) []flight.Flight {
	log.Printf("Query radarbox for zone %s", zone.Name)
	body := r.requestExecutor.Get(getFlightsUrl(zone.Area), headers())
	return parseFlightsResponse(body)
}

func (r RadarboxProvider) GetLink(flight flight.Flight) string {
	return fmt.Sprintf("https://www.radarbox.com/flight/%s", flight.CallSign)
}

func (r RadarboxProvider) MoreDetails(flight flight.Flight) flight.Flight {
	body := r.requestExecutor.Get(getMoreDetailsUrl(flight), headers())
	return parseMoreDetailsResponse(body)
}

func parseMoreDetailsResponse(body string) flight.Flight {
	return flight.Flight{
		Id:           gjson.Get(body, "fid").String(),
		CallSign:     gjson.Get(body, "cs").String(),
		AircraftType: gjson.Get(body, "acd").String(),
		Airline:      gjson.Get(body, "alna").String(),
	}
}
