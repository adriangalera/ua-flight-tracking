package zone

type Zone struct {
	Name string
	Area Bounds
}

type Bounds struct {
	Lat1 float64
	Lng1 float64
	Lat2 float64
	Lng2 float64
}