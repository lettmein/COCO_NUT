package routing

import (
	"context"
	"matcher/internal/geo"
)

type HaversineRouter struct {
	AvgSpeedKmh float64
}

func (h HaversineRouter) Name() string { return "haversine" }

func (h HaversineRouter) Matrix(ctx context.Context, coords []Coord) ([][]float64, error) {
	n := len(coords)
	mat := make([][]float64, n)
	for i := 0; i < n; i++ {
		mat[i] = make([]float64, n)
	}
	if h.AvgSpeedKmh <= 0 {
		h.AvgSpeedKmh = 60
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			a := coords[i]
			b := coords[j]
			km := geo.HaversineKm(geo.Point{Lat: a.Lat, Lon: a.Lon}, geo.Point{Lat: b.Lat, Lon: b.Lon})
			sec := (km / h.AvgSpeedKmh) * 3600.0
			mat[i][j] = sec
			mat[j][i] = sec
		}
	}
	return mat, nil
}
