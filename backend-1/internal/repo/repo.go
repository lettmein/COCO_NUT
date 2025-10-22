package repo

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"matcher/internal/models"
)

type Repo struct {
	DB *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repo { return &Repo{DB: db} }

func (r *Repo) Begin(ctx context.Context) (pgx.Tx, error) {
	return r.DB.Begin(ctx)
}

func (r *Repo) GetRoute(ctx context.Context, tx pgx.Tx, routeID int64) (models.Route, error) {
	var rt models.Route
	err := tx.QueryRow(ctx,
		`SELECT id, depart_at, max_weight_kg, max_volume_m3, status
		   FROM routes WHERE id=$1`, routeID).
		Scan(&rt.ID, &rt.DepartAt, &rt.MaxW, &rt.MaxV, &rt.Status)
	if err != nil {
		return rt, err
	}

	rows, err := tx.Query(ctx,
		`SELECT rp.point_id, rp.seq_no, lp.id, lp.name, lp.lat, lp.lon
		   FROM route_points rp
		   JOIN logistic_points lp ON lp.id = rp.point_id
		  WHERE rp.route_id=$1
		  ORDER BY rp.seq_no ASC`, routeID)
	if err != nil {
		return rt, err
	}
	defer rows.Close()
	for rows.Next() {
		var rp models.RoutePoint
		var p models.Point
		if err := rows.Scan(&rp.PointID, &rp.Seq, &p.ID, &p.Name, &p.Lat, &p.Lon); err != nil {
			return rt, err
		}
		rp.Point = p
		rt.Points = append(rt.Points, rp)
	}
	if len(rt.Points) < 2 {
		return rt, errors.New("route must have at least 2 points")
	}
	return rt, rows.Err()
}

func (r *Repo) LoadPendingRequests(ctx context.Context, tx pgx.Tx, depart time.Time) ([]models.Request, error) {
	rows, err := tx.Query(ctx,
		`SELECT rq.id, rq.origin_point_id, rq.dest_point_id, lp.lat, lp.lon,
		        rq.weight_kg, rq.volume_m3, rq.ready_at, rq.deadline_at, rq.status
		   FROM requests rq
		   JOIN logistic_points lp ON lp.id = rq.dest_point_id
		  WHERE rq.status='pending' AND rq.ready_at <= $1`, depart)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Request
	for rows.Next() {
		var m models.Request
		if err := rows.Scan(&m.ID, &m.OriginID, &m.DestID, &m.Lat, &m.Lon,
			&m.Weight, &m.Volume, &m.ReadyAt, &m.Deadline, &m.Status); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *Repo) SaveAssignments(ctx context.Context, tx pgx.Tx, routeID int64, as []models.Assignment) error {
	// очистить старые
	if _, err := tx.Exec(ctx, `DELETE FROM routes_lists WHERE route_id=$1`, routeID); err != nil {
		return err
	}
	// вставить новые батчем
	batch := &pgx.Batch{}
	for _, a := range as {
		batch.Queue(`INSERT INTO routes_lists(route_id, request_id, seq_no, eta_plan)
		             VALUES($1,$2,$3,$4)
		             ON CONFLICT (route_id, request_id)
		             DO UPDATE SET seq_no=EXCLUDED.seq_no, eta_plan=EXCLUDED.eta_plan`,
			routeID, a.RequestID, a.Seq, a.ETA)
	}
	br := tx.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		return err
	}
	// статусы заявок
	if len(as) > 0 {
		ids := make([]int64, 0, len(as))
		for _, a := range as {
			ids = append(ids, a.RequestID)
		}
		_, _ = tx.Exec(ctx, `UPDATE requests SET status='assigned' WHERE id = ANY($1)`, ids)
	}
	return nil
}

func (r *Repo) GetAssignments(ctx context.Context, routeID int64) ([]models.Assignment, error) {
	rows, err := r.DB.Query(ctx, `SELECT request_id, seq_no, eta_plan FROM routes_lists WHERE route_id=$1 ORDER BY seq_no ASC`, routeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Assignment
	for rows.Next() {
		var a models.Assignment
		if err := rows.Scan(&a.RequestID, &a.Seq, &a.ETA); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *Repo) AdvisoryLock(ctx context.Context, tx pgx.Tx, key int64) (bool, error) {
	var ok bool
	if err := tx.QueryRow(ctx, `SELECT pg_try_advisory_lock($1)`, key).Scan(&ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (r *Repo) FindRoutesToMatch(ctx context.Context, within time.Duration) ([]int64, error) {
	threshold := time.Now().Add(within)
	rows, err := r.DB.Query(ctx,
		`SELECT id
		   FROM routes
		  WHERE status IN ('planned','dispatching')
		    AND depart_at <= $1`, threshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}
