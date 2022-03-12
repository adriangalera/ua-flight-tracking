package flight

type Flight struct {
	Id           string
	CallSign     string
	AircraftType string
	Airline      string
	Link         string
}

type CallSignSet map[string]struct{}

func (s CallSignSet) Add(callSign string) {
	s[callSign] = struct{}{}
}

func (s CallSignSet) Remove(callSign string) {
	delete(s, callSign)
}

func (s CallSignSet) Has(callSign string) bool {
	_, ok := s[callSign]
	return ok
}
