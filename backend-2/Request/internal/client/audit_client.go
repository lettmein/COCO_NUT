package client

import (
	"bytes"
	"encoding/json"
	"git.a7ru.app/a7hack/coco-nut/backend-2/request/internal/dto"
	"net/http"
	"time"
)

type AuditClient struct {
	baseURL string
	client  *http.Client
}

func NewAuditClient(baseURL string) *AuditClient {
	return &AuditClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (ac *AuditClient) SendLog(log *dto.AuditLogDTO) error {
	data, err := json.Marshal(log)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", ac.baseURL+"/api/v1/audit/log", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := ac.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
