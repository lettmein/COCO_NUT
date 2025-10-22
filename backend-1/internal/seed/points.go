package seed

import (
	"context"
	"encoding/json"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LP struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Address string  `json:"address"`
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
}

func SeedLogisticPoints(ctx context.Context, db *pgxpool.Pool, path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var pts []LP
	if err := json.Unmarshal(b, &pts); err != nil {
		return err
	}

	batch := &pgx.Batch{}
	for _, p := range pts {
		batch.Queue(
			`INSERT INTO logistic_points(id, name, address, lat, lon)
			 VALUES($1,$2,$3,$4,$5)
			 ON CONFLICT (id) DO UPDATE SET name=EXCLUDED.name, address=EXCLUDED.address, lat=EXCLUDED.lat, lon=EXCLUDED.lon`,
			p.ID, p.Name, p.Address, p.X, p.Y,
		)
	}
	br := db.SendBatch(ctx, batch)
	return br.Close()
}
