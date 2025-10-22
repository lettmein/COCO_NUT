package dto

import "time"

type AuditLogDTO struct {
	ReqServiceType  string    `json:"req_service_type"`
	RespServiceType string    `json:"resp_service_type"`
	Uri             string    `json:"uri"`
	CreatedAt       time.Time `json:"created_at"`
	DurationTime    float64   `json:"duration_time"`
	RequestBody     string    `json:"request_body"`
	ResponseBody    string    `json:"response_body"`
}
