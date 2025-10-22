package geo

import "math"

type Point struct{ Lat, Lon float64 }

const earthR = 6371.0088

func rad(v float64) float64 { return v * math.Pi / 180 }

func HaversineKm(a, b Point) float64 {
	dLat := rad(b.Lat - a.Lat)
	dLon := rad(b.Lon - a.Lon)
	lat1 := rad(a.Lat)
	lat2 := rad(b.Lat)
	sinDLat := math.Sin(dLat / 2)
	sinDLon := math.Sin(dLon / 2)
	h := sinDLat*sinDLat + math.Cos(lat1)*math.Cos(lat2)*sinDLon*sinDLon
	return 2 * earthR * math.Asin(math.Sqrt(h))
}
