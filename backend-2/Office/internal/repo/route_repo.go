package repo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"git.a7ru.app/a7hack/coco-nut/backend-2/office/internal/dto"
	"github.com/lib/pq"
	"math"
	"time"
)

type RouteRepository struct {
	db *sql.DB
}

func NewRouteRepository(db *sql.DB) *RouteRepository {
	return &RouteRepository{db: db}
}

func (r *RouteRepository) Create(route *dto.CreateRouteDTO) (*dto.RouteResponse, error) {
	now := time.Now()

	routePointsWithTime := r.calculateArrivalTimes(route.RoutePoints, route.DepartureDate)
	routePointsJSON, err := json.Marshal(routePointsWithTime)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO routes (
			depart_at, max_weight_kg, max_volume_m3,
			current_volume, current_weight,
			route_points, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	var id int64
	var createdAt, updatedAt time.Time
	err = r.db.QueryRow(
		query,
		route.DepartureDate, route.MaxWeight, route.MaxVolume,
		0.0, 0.0,
		routePointsJSON, "pending", now, now,
	).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		return nil, err
	}

	if err := r.createRoutePoints(id, route.RoutePoints); err != nil {
		return nil, err
	}

	return &dto.RouteResponse{
		ID:            fmt.Sprintf("%d", id),
		MaxVolume:     route.MaxVolume,
		MaxWeight:     route.MaxWeight,
		CurrentVolume: 0.0,
		CurrentWeight: 0.0,
		DepartureDate: route.DepartureDate,
		RoutePoints:   routePointsWithTime,
		Status:        "pending",
		RequestIDs:    []string{},
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}, nil
}

func (r *RouteRepository) GetByID(id string) (*dto.RouteResponse, error) {
	query := `
		SELECT
			id, max_volume_m3, max_weight_kg, current_volume, current_weight,
			depart_at, route_points, status, request_ids, created_at, updated_at
		FROM routes WHERE id = $1
	`

	var resp dto.RouteResponse
	var routePointsJSON []byte
	var requestIDs pq.StringArray
	var dbID int64

	err := r.db.QueryRow(query, id).Scan(
		&dbID, &resp.MaxVolume, &resp.MaxWeight, &resp.CurrentVolume, &resp.CurrentWeight,
		&resp.DepartureDate, &routePointsJSON, &resp.Status, &requestIDs,
		&resp.CreatedAt, &resp.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(routePointsJSON, &resp.RoutePoints); err != nil {
		return nil, err
	}

	resp.ID = fmt.Sprintf("%d", dbID)
	resp.RequestIDs = requestIDs

	return &resp, nil
}

func (r *RouteRepository) GetAll() ([]*dto.RouteResponse, error) {
	query := `
		SELECT
			id, max_volume_m3, max_weight_kg, current_volume, current_weight,
			depart_at, route_points, status, request_ids, created_at, updated_at
		FROM routes ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes []*dto.RouteResponse
	for rows.Next() {
		var route dto.RouteResponse
		var routePointsJSON []byte
		var requestIDs pq.StringArray
		var dbID int64

		err := rows.Scan(
			&dbID, &route.MaxVolume, &route.MaxWeight, &route.CurrentVolume, &route.CurrentWeight,
			&route.DepartureDate, &routePointsJSON, &route.Status, &requestIDs,
			&route.CreatedAt, &route.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(routePointsJSON, &route.RoutePoints); err != nil {
			return nil, err
		}

		route.ID = fmt.Sprintf("%d", dbID)
		route.RequestIDs = requestIDs
		routes = append(routes, &route)
	}

	return routes, nil
}

func (r *RouteRepository) UpdateStatus(id string, status string) error {
	query := `UPDATE routes SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

func (r *RouteRepository) Delete(id string) error {
	query := `DELETE FROM routes WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *RouteRepository) AddRequestToRoute(routeID string, requestID string, weight float64, volume float64) error {
	query := `
		UPDATE routes
		SET request_ids = array_append(request_ids, $1),
		    current_weight = current_weight + $2,
		    current_volume = current_volume + $3,
		    updated_at = $4
		WHERE id = $5
	`
	_, err := r.db.Exec(query, requestID, weight, volume, time.Now(), routeID)
	return err
}

func (r *RouteRepository) GetByStatus(status string) ([]*dto.RouteResponse, error) {
	query := `
		SELECT
			id, max_volume_m3, max_weight_kg, current_volume, current_weight,
			depart_at, route_points, status, request_ids, created_at, updated_at
		FROM routes WHERE status = $1 ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes []*dto.RouteResponse
	for rows.Next() {
		var route dto.RouteResponse
		var routePointsJSON []byte
		var requestIDs pq.StringArray
		var dbID int64

		err := rows.Scan(
			&dbID, &route.MaxVolume, &route.MaxWeight, &route.CurrentVolume, &route.CurrentWeight,
			&route.DepartureDate, &routePointsJSON, &route.Status, &requestIDs,
			&route.CreatedAt, &route.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(routePointsJSON, &route.RoutePoints); err != nil {
			return nil, err
		}

		route.ID = fmt.Sprintf("%d", dbID)
		route.RequestIDs = requestIDs
		routes = append(routes, &route)
	}

	return routes, nil
}

func (r *RouteRepository) calculateArrivalTimes(points []dto.RoutePoint, departureDate time.Time) []dto.RoutePoint {
	if len(points) == 0 {
		return points
	}

	result := make([]dto.RoutePoint, len(points))
	currentTime := departureDate

	for i := 0; i < len(points); i++ {
		result[i] = points[i]
		result[i].ArrivalTime = currentTime

		if i < len(points)-1 {
			distance := haversineDistance(
				points[i].Latitude, points[i].Longitude,
				points[i+1].Latitude, points[i+1].Longitude,
			)
			hours := distance / 60.0
			currentTime = currentTime.Add(time.Duration(hours * float64(time.Hour)))
		}
	}

	return result
}

func (r *RouteRepository) createRoutePoints(routeID int64, points []dto.RoutePoint) error {
	for i, point := range points {
		pointID, err := r.getOrCreateLogisticPoint(point.Address, point.Latitude, point.Longitude)
		if err != nil {
			return err
		}

		query := `
			INSERT INTO route_points (route_id, seq_no, point_id)
			VALUES ($1, $2, $3)
		`
		_, err = r.db.Exec(query, routeID, i, pointID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RouteRepository) getOrCreateLogisticPoint(address string, lat, lon float64) (int64, error) {
	var id int64
	query := `SELECT id FROM logistic_points WHERE lat = $1 AND lon = $2 LIMIT 1`
	err := r.db.QueryRow(query, lat, lon).Scan(&id)
	if err == nil {
		return id, nil
	}

	if err == sql.ErrNoRows {
		insertQuery := `
			INSERT INTO logistic_points (name, address, lat, lon)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`
		err = r.db.QueryRow(insertQuery, address, address, lat, lon).Scan(&id)
		return id, err
	}

	return 0, err
}

func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371.0

	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degreesToRadians(lat1))*math.Cos(degreesToRadians(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
