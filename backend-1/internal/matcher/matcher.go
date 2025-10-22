package matcher

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"sort"
	"time"

	"matcher/internal/audit"
	"matcher/internal/models"
	"matcher/internal/repo"
	"matcher/internal/routing"
)

type Service struct {
	Repo          *repo.Repo
	Router        routing.RouterClient
	Audit         *audit.Client
	MaxDetourSec  float64
	RouteRadiusKm float64
	AvgSpeedKmh   float64
}

func NewService(r *repo.Repo, rc routing.RouterClient, aud *audit.Client, maxDetourMin float64, radiusKm float64, avgSpeed float64) *Service {
	return &Service{
		Repo:          r,
		Router:        rc,
		Audit:         aud,
		MaxDetourSec:  maxDetourMin * 60.0,
		RouteRadiusKm: radiusKm,
		AvgSpeedKmh:   avgSpeed,
	}
}

// MatchRoute подбирает заявки и записывает их в routes_lists. Никогда не возвращает nil, только [].
func (s *Service) MatchRoute(ctx context.Context, routeID int64) ([]models.Assignment, error) {
	tx, err := s.Repo.Begin(ctx)
	if err != nil {
		return []models.Assignment{}, err
	}
	defer tx.Rollback(ctx)

	ok, err := s.Repo.AdvisoryLock(ctx, tx, routeID)
	if err != nil {
		return []models.Assignment{}, err
	}
	if !ok {
		// кто-то уже считает (например, воркер)
		log.Printf("match route=%d: skipped (lock busy)", routeID)
		return []models.Assignment{}, nil
	}

	route, err := s.Repo.GetRoute(ctx, tx, routeID)
	if err != nil {
		return []models.Assignment{}, err
	}
	if len(route.Points) < 2 {
		log.Printf("match route=%d: invalid route points", routeID)
		return []models.Assignment{}, fmt.Errorf("route %d: need at least 2 points", routeID)
	}

	cands, err := s.Repo.LoadPendingRequests(ctx, tx, route.DepartAt)
	if err != nil {
		return []models.Assignment{}, err
	}
	if len(cands) == 0 {
		log.Printf("match route=%d: no candidates (pending/ready)", routeID)
		_ = tx.Commit(ctx)
		return []models.Assignment{}, nil
	}

	// матрица времени по дорогам: [start, candidates..., end]
	start := route.Points[0].Point
	end := route.Points[len(route.Points)-1].Point

	allCoords := make([]routing.Coord, 0, len(cands)+2)
	allCoords = append(allCoords, routing.Coord{Lat: start.Lat, Lon: start.Lon})
	for _, c := range cands {
		allCoords = append(allCoords, routing.Coord{Lat: c.Lat, Lon: c.Lon})
	}
	allCoords = append(allCoords, routing.Coord{Lat: end.Lat, Lon: end.Lon})

	mat, err := s.Router.Matrix(ctx, allCoords)
	if err != nil {
		return []models.Assignment{}, err
	}

	base := mat[0][len(allCoords)-1]

	// скоринг кандидатов
	type scored struct {
		Req   models.Request
		Score float64
		Det   float64
	}
	scoredC := make([]scored, 0, len(cands))
	for i, r := range cands {
		idx := i + 1
		detour := mat[0][idx] + mat[idx][len(allCoords)-1] - base
		if s.MaxDetourSec > 0 && detour > s.MaxDetourSec {
			continue
		}
		urg := urgency(route.DepartAt, r.Deadline)
		detPenalty := 1.0
		if s.MaxDetourSec > 0 {
			detPenalty = 1.0 - math.Min(1.0, detour/s.MaxDetourSec)
		}
		score := 0.45*1.0 + 0.35*urg + 0.20*detPenalty
		scoredC = append(scoredC, scored{Req: r, Score: score, Det: detour})
	}

	// жадный 2D рюкзак
	sort.Slice(scoredC, func(i, j int) bool { return scoredC[i].Score > scoredC[j].Score })
	wRem, vRem := route.MaxW, route.MaxV
	picked := make([]scored, 0, len(scoredC))
	for _, c := range scoredC {
		wShare := c.Req.Weight / max(1e-9, wRem)
		vShare := c.Req.Volume / max(1e-9, vRem)
		dom := math.Max(wShare, vShare)
		sizeFit := 1.0 - math.Min(1.0, dom)
		c.Score = 0.45*sizeFit + 0.35*urgency(route.DepartAt, c.Req.Deadline) + 0.20*(1.0-math.Min(1.0, c.Det/s.MaxDetourSec))
		if c.Req.Weight <= wRem && c.Req.Volume <= vRem {
			picked = append(picked, c)
			wRem -= c.Req.Weight
			vRem -= c.Req.Volume
		}
	}

	if len(picked) == 0 {
		log.Printf("match route=%d: zero picked (constraints/detour)", routeID)
		_ = tx.Commit(ctx)
		return []models.Assignment{}, nil
	}

	// соберём матрицу для [start + picked + end]
	coords := make([]routing.Coord, 0, len(picked)+2)
	coords = append(coords, routing.Coord{Lat: start.Lat, Lon: start.Lon})
	for _, p := range picked {
		coords = append(coords, routing.Coord{Lat: p.Req.Lat, Lon: p.Req.Lon})
	}
	coords = append(coords, routing.Coord{Lat: end.Lat, Lon: end.Lon})

	mat2, err := s.Router.Matrix(ctx, coords)
	if err != nil {
		return []models.Assignment{}, err
	}

	// порядок: NN + 2-opt, фикс. начало/конец
	order := solvePathFixed(mat2)

	// ETA
	assign := make([]models.Assignment, 0, len(picked))
	cur := route.DepartAt
	for seq := 1; seq < len(order); seq++ {
		from := order[seq-1]
		to := order[seq]
		dur := mat2[from][to]
		cur = cur.Add(time.Duration(dur * float64(time.Second)))
		if to == len(order)-1 { // конец маршрута
			break
		}
		pIndex := to - 1
		assign = append(assign, models.Assignment{
			RequestID: picked[pIndex].Req.ID,
			Seq:       seq,
			ETA:       cur,
		})
	}

	if err := s.Repo.SaveAssignments(ctx, tx, route.ID, assign); err != nil {
		return []models.Assignment{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return []models.Assignment{}, err
	}

	// аудит — best-effort
	if s.Audit != nil {
		s.Audit.Send(ctx, audit.Event{
			FromService: "matcher",
			ToService:   s.Router.Name(),
			URI:         "matrix",
			HTTPStatus:  200,
			At:          time.Now().Format(time.RFC3339),
			DurationMs:  0,
		})
	}

	log.Printf("match route=%d: saved %d assignments", route.ID, len(assign))
	return assign, nil
}

func urgency(depart time.Time, dl sql.NullTime) float64 {
	if !dl.Valid {
		return 0.5
	}
	if dl.Time.Before(depart) {
		return 1.0
	}
	d := dl.Time.Sub(depart).Hours()
	if d <= 0 {
		return 1
	}
	if d > 72 {
		return 0
	}
	return 1 - d/72.0
}

func solvePathFixed(mat [][]float64) []int {
	n := len(mat)
	if n <= 2 {
		order := make([]int, n)
		for i := range order {
			order[i] = i
		}
		return order
	}
	visited := make([]bool, n)
	order := make([]int, 0, n)
	cur := 0
	visited[cur] = true
	order = append(order, cur)
	for len(order) < n-1 {
		best, bestW := -1, math.MaxFloat64
		for j := 1; j < n-1; j++ {
			if visited[j] {
				continue
			}
			if w := mat[cur][j]; w < bestW {
				bestW, best = w, j
			}
		}
		if best == -1 {
			break
		}
		visited[best] = true
		order = append(order, best)
		cur = best
	}
	order = append(order, n-1)
	// 2-opt c фикс. концами
	improved := true
	for improved {
		improved = false
		for i := 1; i < len(order)-3; i++ {
			for j := i + 1; j < len(order)-1; j++ {
				a1, a2 := order[i-1], order[i]
				b1, b2 := order[j], order[j+1]
				old := mat[a1][a2] + mat[b1][b2]
				newW := mat[a1][b1] + mat[a2][b2]
				if newW+1e-6 < old {
					for l, r := i, j; l < r; l, r = l+1, r-1 {
						order[l], order[r] = order[r], order[l]
					}
					improved = true
				}
			}
		}
	}
	return order
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
