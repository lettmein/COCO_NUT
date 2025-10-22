package service

import (
	"git.a7ru.app/a7hack/coco-nut/backend-2/office/internal/dto"
	"git.a7ru.app/a7hack/coco-nut/backend-2/office/internal/repo"
	"math"
	"time"
)

const TRUCK_SPEED_KMH = 60.0

type RouteService struct {
	repo *repo.RouteRepository
}

func NewRouteService(repo *repo.RouteRepository) *RouteService {
	return &RouteService{repo: repo}
}

func (s *RouteService) CreateRoute(route *dto.CreateRouteDTO) (*dto.RouteResponse, error) {
	s.calculateArrivalTimes(route)
	return s.repo.Create(route)
}

func (s *RouteService) GetRoute(id string) (*dto.RouteResponse, error) {
	return s.repo.GetByID(id)
}

func (s *RouteService) GetAllRoutes() ([]*dto.RouteResponse, error) {
	return s.repo.GetAll()
}

func (s *RouteService) UpdateRouteStatus(id string, status string) error {
	return s.repo.UpdateStatus(id, status)
}

func (s *RouteService) DeleteRoute(id string) error {
	return s.repo.Delete(id)
}

func (s *RouteService) GetRoutesByStatus(status string) ([]*dto.RouteResponse, error) {
	return s.repo.GetByStatus(status)
}

func (s *RouteService) AddRequestToRoute(routeID string, requestID string, weight float64, volume float64) error {
	return s.repo.AddRequestToRoute(routeID, requestID, weight, volume)
}

func (s *RouteService) calculateArrivalTimes(route *dto.CreateRouteDTO) {
	if len(route.RoutePoints) == 0 {
		return
	}

	route.RoutePoints[0].ArrivalTime = route.DepartureDate

	for i := 1; i < len(route.RoutePoints); i++ {
		distance := s.calculateDistance(
			route.RoutePoints[i-1].Latitude, route.RoutePoints[i-1].Longitude,
			route.RoutePoints[i].Latitude, route.RoutePoints[i].Longitude,
		)

		travelTimeHours := distance / TRUCK_SPEED_KMH
		travelTimeMinutes := int(travelTimeHours * 60)

		route.RoutePoints[i].ArrivalTime = route.RoutePoints[i-1].ArrivalTime.Add(
			time.Duration(travelTimeMinutes) * time.Minute,
		)
	}
}

func (s *RouteService) calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371.0

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}
