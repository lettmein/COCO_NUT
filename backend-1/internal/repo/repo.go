package repo

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"matcher/internal/models"
)

// Request statuses
const (
	RequestStatusPending   = "pending"
	RequestStatusAssigned  = "assigned"
	RequestStatusInTransit = "in_transit"
	RequestStatusDelivered = "delivered"
	RequestStatusCancelled = "cancelled"
)

// Route statuses
const (
	RouteStatusPlanned    = "planned"
	RouteStatusPending    = "pending"
	RouteStatusInProgress = "in_progress"
	RouteStatusCompleted  = "completed"
	RouteStatusCancelled  = "cancelled"
)

type Repo struct {
	DB *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repo { return &Repo{DB: db} }

func (r *Repo) Begin(ctx context.Context) (pgx.Tx, error) {
	return r.DB.Begin(ctx)
}

// GetRoute: базовая инфа + опорные точки маршрута
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

// LoadCandidates: pending + уже назначенные в ЭТОТ маршрут (чтобы переоценить их заново)
func (r *Repo) LoadCandidates(ctx context.Context, tx pgx.Tx, routeID int64, depart time.Time) ([]models.Request, error) {
	rows, err := tx.Query(ctx,
		`SELECT DISTINCT rq.id, rq.origin_point_id, rq.dest_point_id, lp.lat, lp.lon,
		        rq.weight_kg, rq.volume_m3, rq.ready_at, rq.deadline_at, rq.status
		   FROM requests rq
		   JOIN logistic_points lp ON lp.id = rq.dest_point_id
		   LEFT JOIN routes_lists rl
		     ON rl.request_id = rq.id AND rl.route_id = $1
		  WHERE (rq.status = $2 AND rq.ready_at <= $3)
		     OR rl.route_id IS NOT NULL`,
		routeID, RequestStatusPending, depart)
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

// SaveAssignments: полностью переиздаёт лист маршрута и синхронизирует статусы и агрегаты
//   - удаляет старый лист, вставляет новый (в порядке seq);
//   - новые заявки -> assigned;
//   - выпавшие из листа -> pending (если не назначены в другие маршруты);
//   - обновляет routes.current_weight/current_volume/request_ids.
func (r *Repo) SaveAssignments(ctx context.Context, tx pgx.Tx, routeID int64, as []models.Assignment) error {
	// 1) предыдущие назначения
	prevIDs := make([]int64, 0)
	rows, err := tx.Query(ctx, `SELECT request_id FROM routes_lists WHERE route_id=$1`, routeID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return err
		}
		prevIDs = append(prevIDs, id)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	// 2) перезапись листа
	if _, err := tx.Exec(ctx, `DELETE FROM routes_lists WHERE route_id=$1`, routeID); err != nil {
		return err
	}
	batch := &pgx.Batch{}
	newIDs := make([]int64, 0, len(as))
	newIDsText := make([]string, 0, len(as)) // для routes.request_ids (text[])
	for _, a := range as {
		newIDs = append(newIDs, a.RequestID)
		newIDsText = append(newIDsText, strconv.FormatInt(a.RequestID, 10))
		batch.Queue(`INSERT INTO routes_lists(route_id, request_id, seq_no, eta_plan)
		             VALUES($1,$2,$3,$4)
		             ON CONFLICT (route_id, request_id)
		             DO UPDATE SET seq_no=EXCLUDED.seq_no, eta_plan=EXCLUDED.eta_plan`,
			routeID, a.RequestID, a.Seq, a.ETA)
	}
	if len(as) > 0 {
		br := tx.SendBatch(ctx, batch)
		if err := br.Close(); err != nil {
			return err
		}
	}

	// 3) новые -> assigned
	if len(newIDs) > 0 {
		_, _ = tx.Exec(ctx,
			`UPDATE requests
			    SET status=$1, updated_at=now()
			  WHERE id = ANY($2::bigint[])`,
			RequestStatusAssigned, newIDs)
	}

	// 4) выпавшие -> pending (если больше не назначены в др. маршруты)
	removed := diffIDs(prevIDs, newIDs)
	if len(removed) > 0 {
		_, _ = tx.Exec(ctx,
			`UPDATE requests
			    SET status=$1, updated_at=now()
			  WHERE id = ANY($2::bigint[])
			    AND status=$3
			    AND NOT EXISTS (SELECT 1 FROM routes_lists rl WHERE rl.request_id = requests.id)`,
			RequestStatusPending, removed, RequestStatusAssigned)
	}

	// 5) агрегаты маршрута
	_, _ = tx.Exec(ctx,
		`UPDATE routes
		    SET current_weight = COALESCE((SELECT SUM(weight_kg) FROM requests WHERE id = ANY($1::bigint[])), 0::numeric),
		        current_volume = COALESCE((SELECT SUM(volume_m3) FROM requests WHERE id = ANY($1::bigint[])), 0::numeric),
		        request_ids    = $2::text[],
		        updated_at     = now()
		  WHERE id = $3`,
		newIDs, newIDsText, routeID)

	return nil
}

func (r *Repo) GetAssignments(ctx context.Context, routeID int64) ([]models.Assignment, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT request_id, seq_no, eta_plan
		   FROM routes_lists
		  WHERE route_id=$1
		  ORDER BY seq_no ASC`, routeID)
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

// транзакционный advisory lock
func (r *Repo) AdvisoryLock(ctx context.Context, tx pgx.Tx, key int64) (bool, error) {
	var ok bool
	if err := tx.QueryRow(ctx, `SELECT pg_try_advisory_xact_lock($1)`, key).Scan(&ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (r *Repo) FindRoutesToMatch(ctx context.Context, within time.Duration) ([]int64, error) {
	threshold := time.Now().Add(within)
	rows, err := r.DB.Query(ctx,
		`SELECT id
		   FROM routes
		  WHERE status IN ($1,$2)
		    AND depart_at <= $3
		  ORDER BY depart_at ASC`,
		RouteStatusPlanned, RouteStatusPending, threshold)
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

// helpers

func diffIDs(prev, cur []int64) []int64 {
	set := make(map[int64]struct{}, len(cur))
	for _, id := range cur {
		set[id] = struct{}{}
	}
	out := make([]int64, 0, len(prev))
	for _, id := range prev {
		if _, ok := set[id]; !ok {
			out = append(out, id)
		}
	}
	return out
}
