package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Event struct {
	FromService string  `json:"from_service"`
	ToService   string  `json:"to_service"`
	URI         string  `json:"uri"`
	HTTPStatus  int     `json:"http_status"`
	At          string  `json:"at"`
	DurationMs  int64   `json:"duration_ms"`
	UserID      *string `json:"user_id,omitempty"`
	ReqBody     string  `json:"request_body,omitempty"`
	RespBody    string  `json:"response_body,omitempty"`
}

type Client struct {
	BaseURL string
	Client  *http.Client
	From    string
}

func (c *Client) Send(ctx context.Context, ev Event) {
	if c == nil || c.BaseURL == "" {
		return
	}
	if c.Client == nil {
		c.Client = &http.Client{Timeout: 5 * time.Second}
	}
	b, _ := json.Marshal(ev)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/audit/events", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Client.Do(req)
	if err != nil {
		return
	}
	resp.Body.Close()
}
