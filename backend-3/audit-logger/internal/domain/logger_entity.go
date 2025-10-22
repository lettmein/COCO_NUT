package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidServiceType = errors.New("invalid service type")
	ErrInvalidURI         = errors.New("invalid URI")
	ErrInvalidID          = errors.New("invalid ID")
	ErrInvalidTimeRange   = errors.New("invalid time range")
	ErrLogNotFound        = errors.New("log not found")
)

type Log struct {
	ID              int       `json:"id"`
	ReqServiceType  string    `json:"req_service_type"`
	RespServiceType string    `json:"resp_service_type"`
	Uri             string    `json:"uri"`
	CreatedAt       time.Time `json:"created_at"`
	DurationTime    float64   `json:"duration_time"`
	RequestBody     string    `json:"request_body"`
	ResponseBody    string    `json:"response_body"`
}

type LoggerRepository interface {
	Create(ctx context.Context, log *Log) error
	GetWithFilter(ctx context.Context, filter LogFilter) ([]Log, error)
}

type LoggerService interface {
	Create(ctx context.Context, log *Log) error
	GetAllWithFilter(ctx context.Context, filter LogFilter) ([]Log, error)
}

type LogFilter struct {
	ServiceType string
	URI         string
	StartTime   time.Time
	EndTime     time.Time
	MinDuration float64
	MaxDuration float64
	Limit       int
	Offset      int
}
