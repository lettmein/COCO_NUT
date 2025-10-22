package repo

import (
	"database/sql"
	"encoding/json"
	"git.a7ru.app/a7hack/coco-nut/backend-2/office/internal/dto"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

type RouteRepository struct {
	db *sql.DB
}

func NewRouteRepository(db *sql.DB) *RouteRepository {
	return &RouteRepository{db: db}
}

func (r *RouteRepository) Create(route *dto.CreateRouteDTO) (*dto.RouteResponse, error) {
	id := uuid.New().String()
	now := time.Now()

	routePointsJSON, err := json.Marshal(route.RoutePoints)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO routes (
			id, max_volume, max_weight, current_volume, current_weight,
			departure_date, route_points, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`

	var createdAt, updatedAt time.Time
	err = r.db.QueryRow(
		query,
		id, route.MaxVolume, route.MaxWeight, 0.0, 0.0,
		route.DepartureDate, routePointsJSON, "pending", now, now,
	).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		return nil, err
	}

	return &dto.RouteResponse{
		ID:            id,
		MaxVolume:     route.MaxVolume,
		MaxWeight:     route.MaxWeight,
		CurrentVolume: 0.0,
		CurrentWeight: 0.0,
		DepartureDate: route.DepartureDate,
		RoutePoints:   route.RoutePoints,
		Status:        "pending",
		RequestIDs:    []string{},
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}, nil
}

func (r *RouteRepository) GetByID(id string) (*dto.RouteResponse, error) {
	query := `
		SELECT
			id, max_volume, max_weight, current_volume, current_weight,
			departure_date, route_points, status, request_ids, created_at, updated_at
		FROM routes WHERE id = $1
	`

	var resp dto.RouteResponse
	var routePointsJSON []byte
	var requestIDs pq.StringArray

	err := r.db.QueryRow(query, id).Scan(
		&resp.ID, &resp.MaxVolume, &resp.MaxWeight, &resp.CurrentVolume, &resp.CurrentWeight,
		&resp.DepartureDate, &routePointsJSON, &resp.Status, &requestIDs,
		&resp.CreatedAt, &resp.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(routePointsJSON, &resp.RoutePoints); err != nil {
		return nil, err
	}

	resp.RequestIDs = requestIDs

	return &resp, nil
}

func (r *RouteRepository) GetAll() ([]*dto.RouteResponse, error) {
	query := `
		SELECT
			id, max_volume, max_weight, current_volume, current_weight,
			departure_date, route_points, status, request_ids, created_at, updated_at
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

		err := rows.Scan(
			&route.ID, &route.MaxVolume, &route.MaxWeight, &route.CurrentVolume, &route.CurrentWeight,
			&route.DepartureDate, &routePointsJSON, &route.Status, &requestIDs,
			&route.CreatedAt, &route.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(routePointsJSON, &route.RoutePoints); err != nil {
			return nil, err
		}

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
			id, max_volume, max_weight, current_volume, current_weight,
			departure_date, route_points, status, request_ids, created_at, updated_at
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

		err := rows.Scan(
			&route.ID, &route.MaxVolume, &route.MaxWeight, &route.CurrentVolume, &route.CurrentWeight,
			&route.DepartureDate, &routePointsJSON, &route.Status, &requestIDs,
			&route.CreatedAt, &route.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(routePointsJSON, &route.RoutePoints); err != nil {
			return nil, err
		}

		route.RequestIDs = requestIDs
		routes = append(routes, &route)
	}

	return routes, nil
}
