package dto

import "time"

type RoutePoint struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Address   string    `json:"address"`
	ArrivalTime time.Time `json:"arrival_time,omitempty"`
}

type CreateRouteDTO struct {
	MaxVolume    float64      `json:"max_volume"`
	MaxWeight    float64      `json:"max_weight"`
	DepartureDate time.Time   `json:"departure_date"`
	RoutePoints  []RoutePoint `json:"route_points"`
}

type UpdateRouteDTO struct {
	Status string `json:"status"`
}

type RouteResponse struct {
	ID            string       `json:"id"`
	MaxVolume     float64      `json:"max_volume"`
	MaxWeight     float64      `json:"max_weight"`
	CurrentVolume float64      `json:"current_volume"`
	CurrentWeight float64      `json:"current_weight"`
	DepartureDate time.Time    `json:"departure_date"`
	RoutePoints   []RoutePoint `json:"route_points"`
	Status        string       `json:"status"`
	RequestIDs    []string     `json:"request_ids,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
