package repository

import (
	"context"
	"database/sql"
	"fmt"
	"myapp/internal/domain"
)

type loggerRepository struct {
	db *sql.DB
}

func NewLoggerRepository(db *sql.DB) domain.LoggerRepository {
	return &loggerRepository{db: db}
}

func (r *loggerRepository) Create(ctx context.Context, log *domain.Log) error {
	query := `
		INSERT INTO logs 
		(req_service_type, resp_service_type, uri, created_at, duration_time, request_body, response_body) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query,
		log.ReqServiceType,
		log.RespServiceType,
		log.Uri,
		log.CreatedAt,
		log.DurationTime,
		log.RequestBody,
		log.ResponseBody,
	).Scan(&log.ID)

	return err
}

func (r *loggerRepository) GetByID(ctx context.Context, id int) (*domain.Log, error) {
	query := `
		SELECT id, req_service_type, resp_service_type, uri, created_at, duration_time, request_body, response_body
		FROM logs WHERE id = $1
	`

	log := &domain.Log{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&log.ID,
		&log.ReqServiceType,
		&log.RespServiceType,
		&log.Uri,
		&log.CreatedAt,
		&log.DurationTime,
		&log.RequestBody,
		&log.ResponseBody,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return log, err
}

func (r *loggerRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.Log, error) {
	query := `
		SELECT id, req_service_type, resp_service_type, uri, created_at, duration_time, request_body, response_body
		FROM logs 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.Log
	for rows.Next() {
		log := &domain.Log{}
		err := rows.Scan(
			&log.ID,
			&log.ReqServiceType,
			&log.RespServiceType,
			&log.Uri,
			&log.CreatedAt,
			&log.DurationTime,
			&log.RequestBody,
			&log.ResponseBody,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func (r *loggerRepository) Update(ctx context.Context, id int, log *domain.Log) error {
	query := `
		UPDATE logs 
		SET req_service_type = $1, resp_service_type = $2, uri = $3, duration_time = $4, 
		    request_body = $5, response_body = $6
		WHERE id = $7
	`

	result, err := r.db.ExecContext(ctx, query,
		log.ReqServiceType,
		log.RespServiceType,
		log.Uri,
		log.DurationTime,
		log.RequestBody,
		log.ResponseBody,
		id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *loggerRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM logs WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Дополнительный метод для фильтрации (может пригодиться для сервиса)
func (r *loggerRepository) GetWithFilter(ctx context.Context, filter domain.LogFilter) ([]*domain.Log, error) {
	query := `SELECT id, req_service_type, resp_service_type, uri, created_at, duration_time, request_body, response_body FROM logs WHERE 1=1`
	args := []interface{}{}
	argPos := 1

	if filter.ServiceType != "" {
		query += fmt.Sprintf(" AND (req_service_type = $%d OR resp_service_type = $%d)", argPos, argPos)
		args = append(args, filter.ServiceType)
		argPos++
	}

	if filter.URI != "" {
		query += fmt.Sprintf(" AND uri LIKE $%d", argPos)
		args = append(args, "%"+filter.URI+"%")
		argPos++
	}

	if !filter.StartTime.IsZero() {
		query += fmt.Sprintf(" AND created_at >= $%d", argPos)
		args = append(args, filter.StartTime)
		argPos++
	}

	if !filter.EndTime.IsZero() {
		query += fmt.Sprintf(" AND created_at <= $%d", argPos)
		args = append(args, filter.EndTime)
		argPos++
	}

	if filter.MinDuration > 0 {
		query += fmt.Sprintf(" AND duration_time >= $%d", argPos)
		args = append(args, filter.MinDuration)
		argPos++
	}

	if filter.MaxDuration > 0 {
		query += fmt.Sprintf(" AND duration_time <= $%d", argPos)
		args = append(args, filter.MaxDuration)
		argPos++
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, filter.Limit)
		argPos++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argPos)
		args = append(args, filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.Log
	for rows.Next() {
		log := &domain.Log{}
		err := rows.Scan(
			&log.ID,
			&log.ReqServiceType,
			&log.RespServiceType,
			&log.Uri,
			&log.CreatedAt,
			&log.DurationTime,
			&log.RequestBody,
			&log.ResponseBody,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}
