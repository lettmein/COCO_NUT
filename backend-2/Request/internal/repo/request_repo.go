package repo

import (
	"database/sql"
	"fmt"
	"git.a7ru.app/a7hack/coco-nut/backend-2/request/internal/dto"
	"time"
)

type RequestRepository struct {
	db *sql.DB
}

func NewRequestRepository(db *sql.DB) *RequestRepository {
	return &RequestRepository{db: db}
}

func (r *RequestRepository) Create(req *dto.CreateRequestDTO) (*dto.RequestResponse, error) {
	now := time.Now()

	originPointID, err := r.getOrCreateLogisticPoint(req.LogisticPointID)
	if err != nil {
		return nil, err
	}

	destPointID, err := r.createDestinationPoint(req.Recipient.Address, 0, 0)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO requests (
			origin_point_id, dest_point_id,
			weight_kg, volume_m3, ready_at, deadline_at,
			customer_company_name, customer_inn, customer_contact_name, customer_phone, customer_email,
			cargo_name, cargo_quantity, cargo_special_requirements,
			recipient_company_name, recipient_address, recipient_contact_name, recipient_phone,
			status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
		RETURNING id, created_at, updated_at
	`

	var id int64
	var createdAt, updatedAt time.Time
	err = r.db.QueryRow(
		query,
		originPointID, destPointID,
		req.Cargo.Weight, req.Cargo.Volume, now, nil,
		req.Customer.CompanyName, req.Customer.INN, req.Customer.ContactName, req.Customer.Phone, req.Customer.Email,
		req.Cargo.Name, req.Cargo.Quantity, req.Cargo.SpecialRequirements,
		req.Recipient.CompanyName, req.Recipient.Address, req.Recipient.ContactName, req.Recipient.Phone,
		"pending", now, now,
	).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		return nil, err
	}

	return &dto.RequestResponse{
		ID:              fmt.Sprintf("%d", id),
		LogisticPointID: req.LogisticPointID,
		Customer:        req.Customer,
		Cargo:           req.Cargo,
		Recipient:       req.Recipient,
		Status:          "pending",
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}, nil
}

func (r *RequestRepository) GetByID(id string) (*dto.RequestResponse, error) {
	query := `
		SELECT
			id, origin_point_id,
			customer_company_name, customer_inn, customer_contact_name, customer_phone, customer_email,
			cargo_name, cargo_quantity, cargo_weight, cargo_volume, cargo_special_requirements,
			recipient_company_name, recipient_address, recipient_contact_name, recipient_phone,
			status, created_at, updated_at
		FROM requests WHERE id = $1
	`

	var resp dto.RequestResponse
	var dbID int64
	var originPointID int64
	err := r.db.QueryRow(query, id).Scan(
		&dbID, &originPointID,
		&resp.Customer.CompanyName, &resp.Customer.INN, &resp.Customer.ContactName, &resp.Customer.Phone, &resp.Customer.Email,
		&resp.Cargo.Name, &resp.Cargo.Quantity, &resp.Cargo.Weight, &resp.Cargo.Volume, &resp.Cargo.SpecialRequirements,
		&resp.Recipient.CompanyName, &resp.Recipient.Address, &resp.Recipient.ContactName, &resp.Recipient.Phone,
		&resp.Status, &resp.CreatedAt, &resp.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	resp.ID = fmt.Sprintf("%d", dbID)
	resp.LogisticPointID = int(originPointID)

	return &resp, nil
}

func (r *RequestRepository) GetAll() ([]*dto.RequestResponse, error) {
	query := `
		SELECT
			id, origin_point_id,
			customer_company_name, customer_inn, customer_contact_name, customer_phone, customer_email,
			cargo_name, cargo_quantity, weight_kg, volume_m3, cargo_special_requirements,
			recipient_company_name, recipient_address, recipient_contact_name, recipient_phone,
			status, created_at, updated_at
		FROM requests ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*dto.RequestResponse
	for rows.Next() {
		var req dto.RequestResponse
		var dbID int64
		var originPointID int64
		err := rows.Scan(
			&dbID, &originPointID,
			&req.Customer.CompanyName, &req.Customer.INN, &req.Customer.ContactName, &req.Customer.Phone, &req.Customer.Email,
			&req.Cargo.Name, &req.Cargo.Quantity, &req.Cargo.Weight, &req.Cargo.Volume, &req.Cargo.SpecialRequirements,
			&req.Recipient.CompanyName, &req.Recipient.Address, &req.Recipient.ContactName, &req.Recipient.Phone,
			&req.Status, &req.CreatedAt, &req.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		req.ID = fmt.Sprintf("%d", dbID)
		req.LogisticPointID = int(originPointID)
		requests = append(requests, &req)
	}

	return requests, nil
}

func (r *RequestRepository) UpdateStatus(id string, status string) error {
	query := `UPDATE requests SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

func (r *RequestRepository) Delete(id string) error {
	query := `DELETE FROM requests WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *RequestRepository) GetByStatus(status string) ([]*dto.RequestResponse, error) {
	query := `
		SELECT
			id, origin_point_id,
			customer_company_name, customer_inn, customer_contact_name, customer_phone, customer_email,
			cargo_name, cargo_quantity, weight_kg, volume_m3, cargo_special_requirements,
			recipient_company_name, recipient_address, recipient_contact_name, recipient_phone,
			status, created_at, updated_at
		FROM requests WHERE status = $1 ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*dto.RequestResponse
	for rows.Next() {
		var req dto.RequestResponse
		var dbID int64
		var originPointID int64
		err := rows.Scan(
			&dbID, &originPointID,
			&req.Customer.CompanyName, &req.Customer.INN, &req.Customer.ContactName, &req.Customer.Phone, &req.Customer.Email,
			&req.Cargo.Name, &req.Cargo.Quantity, &req.Cargo.Weight, &req.Cargo.Volume, &req.Cargo.SpecialRequirements,
			&req.Recipient.CompanyName, &req.Recipient.Address, &req.Recipient.ContactName, &req.Recipient.Phone,
			&req.Status, &req.CreatedAt, &req.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		req.ID = fmt.Sprintf("%d", dbID)
		req.LogisticPointID = int(originPointID)
		requests = append(requests, &req)
	}

	return requests, nil
}

func (r *RequestRepository) getOrCreateLogisticPoint(pointID int) (int64, error) {
	var id int64
	err := r.db.QueryRow("SELECT id FROM logistic_points WHERE id = $1", pointID).Scan(&id)
	if err == nil {
		return id, nil
	}

	if err == sql.ErrNoRows {
		return int64(pointID), nil
	}

	return 0, err
}

func (r *RequestRepository) createDestinationPoint(address string, lat, lon float64) (int64, error) {
	var id int64
	query := `
		INSERT INTO logistic_points (name, address, lat, lon)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err := r.db.QueryRow(query, address, address, lat, lon).Scan(&id)
	return id, err
}
