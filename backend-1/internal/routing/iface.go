package routing

import "context"

type Coord struct {
	Lat float64
	Lon float64
}

type RouterClient interface {
	Name() string
	// Matrix возвращает матрицу времени в секундах для coords x coords
	Matrix(ctx context.Context, coords []Coord) ([][]float64, error)
}
