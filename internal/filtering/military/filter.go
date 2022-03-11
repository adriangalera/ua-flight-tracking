package military

import (
	"agalera.eu/flight-tracking/internal/filtering"
	"agalera.eu/flight-tracking/internal/flight"
	"bufio"
	"log"
	"os"
	"strings"
)

type MilitaryFlightFilter struct {
	knownMilitaryCodes []string
	knownCivilCodes    []string
}

func NewMilitaryFlightFilter(dataFolder string) filtering.Filter {
	return MilitaryFlightFilter{
		knownMilitaryCodes: getKnownMilitaryCodes(dataFolder),
		knownCivilCodes:    getKnownCivilCodes(dataFolder),
	}
}

func getKnownMilitaryCodes(dataFolder string) []string {
	return readAllLinesInFile(dataFolder + "/known-mil-code.txt")
}
func getKnownCivilCodes(dataFolder string) []string {
	return readAllLinesInFile(dataFolder + "/known-civil-code.txt")
}

func readAllLinesInFile(filename string) []string {
	readFile, errorOpening := os.Open(filename)
	if errorOpening != nil {
		log.Fatalf("Could not read the known military codes file! Error: %v", errorOpening)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	errClosing := readFile.Close()
	if errClosing != nil {
		log.Fatalf("Could not close the known military codes file! Error: %v", errClosing)
		return []string{}
	}

	return fileLines
}

func (m MilitaryFlightFilter) isMilitary(flight flight.Flight) bool {
	var military bool
	for _, militaryCallSign := range m.knownMilitaryCodes {
		if strings.HasPrefix(flight.CallSign, militaryCallSign) {
			military = true
			break
		}
	}
	if military {
		for _, civilCallSign := range m.knownCivilCodes {
			if strings.HasPrefix(flight.CallSign, civilCallSign) {
				military = false
				break
			}
		}
	}
	return military
}

func (m MilitaryFlightFilter) Filter(flights []flight.Flight) []flight.Flight {
	var militaryFlights []flight.Flight
	for _, f := range flights {
		if m.isMilitary(f) {
			militaryFlights = append(militaryFlights, f)
		}
	}
	return militaryFlights
}
