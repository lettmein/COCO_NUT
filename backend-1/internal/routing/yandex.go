package routing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Важно: endpoint и формат могут отличаться в твоём аккаунте Яндекса.
// По умолчанию используем v2/matrix с origins/destinations.
// При необходимости поправь URL/тело в соответствии с документацией.

// JSON модели (минимум полей)
type matrixReq struct {
	Origins      []ymPoint `json:"origins"`
	Destinations []ymPoint `json:"destinations"`
	Mode         string    `json:"mode,omitempty"` // "driving"
}

type ymPoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type matrixResp struct {
	// ожидаем поле "matrix" или "rows" — покроем оба.
	Matrix [][]struct {
		TravelTime float64 `json:"travel_time"` // сек
		// Distance может быть, если пригодится
		Distance float64 `json:"distance"`
	} `json:"matrix"`
	Rows []struct {
		Elements []struct {
			Duration struct {
				Value float64 `json:"value"`
			} `json:"duration"`
			Distance struct {
				Value float64 `json:"value"`
			} `json:"distance"`
		} `json:"elements"`
	} `json:"rows"`
}

type YandexRouter struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
	// размер батча (origins x destinations) — консервативно
	MaxBatch int
}

func (y YandexRouter) Name() string { return "yandex" }

func (y YandexRouter) Matrix(ctx context.Context, coords []Coord) ([][]float64, error) {
	if y.Client == nil {
		y.Client = &http.Client{Timeout: 20 * time.Second}
	}
	n := len(coords)
	out := make([][]float64, n)
	for i := range out {
		out[i] = make([]float64, n)
	}
	if n == 0 {
		return out, nil
	}

	// батчим по квадратам
	chunk := y.MaxBatch
	if chunk <= 0 {
		chunk = 25
	}

	for oi := 0; oi < n; oi += chunk {
		oend := min(n, oi+chunk)
		origins := toYM(coords[oi:oend])

		for di := 0; di < n; di += chunk {
			dend := min(n, di+chunk)
			dests := toYM(coords[di:dend])

			reqBody := matrixReq{
				Origins:      origins,
				Destinations: dests,
				Mode:         "driving",
			}
			b, _ := json.Marshal(reqBody)

			url := y.BaseURL
			// ключ пробуем прокинуть как query (или используй заголовок ниже)
			url = fmt.Sprintf("%s?apikey=%s", url, y.APIKey)

			httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
			httpReq.Header.Set("Content-Type", "application/json")
			// альтернативный способ: httpReq.Header.Set("X-Api-Key", y.APIKey)

			resp, err := y.Client.Do(httpReq)
			if err != nil {
				return nil, err
			}
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			if resp.StatusCode/100 != 2 {
				return nil, fmt.Errorf("yandex matrix http %d: %s", resp.StatusCode, string(body))
			}

			var mr matrixResp
			if err := json.Unmarshal(body, &mr); err != nil {
				return nil, err
			}
			// попытка распарсить известные форматы
			switch {
			case len(mr.Matrix) > 0:
				for i := 0; i < oend-oi; i++ {
					for j := 0; j < dend-di; j++ {
						out[oi+i][di+j] = mr.Matrix[i][j].TravelTime
					}
				}
			case len(mr.Rows) > 0:
				for i := 0; i < oend-oi; i++ {
					for j := 0; j < dend-di; j++ {
						out[oi+i][di+j] = mr.Rows[i].Elements[j].Duration.Value
					}
				}
			default:
				return nil, fmt.Errorf("unknown yandex matrix format")
			}
		}
	}

	return out, nil
}

func toYM(cs []Coord) []ymPoint {
	out := make([]ymPoint, len(cs))
	for i, c := range cs {
		out[i] = ymPoint{Lat: c.Lat, Lon: c.Lon}
	}
	return out
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
