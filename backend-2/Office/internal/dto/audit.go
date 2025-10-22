package dto

type AuditLogDTO struct {
	ServiceName    string `json:"service_name"`
	RequestMethod  string `json:"request_method"`
	RequestURL     string `json:"request_url"`
	RequestBody    string `json:"request_body"`
	ResponseBody   string `json:"response_body"`
	StatusCode     int    `json:"status_code"`
	UserID         string `json:"user_id,omitempty"`
	ProcessingTime int64  `json:"processing_time"`
}
