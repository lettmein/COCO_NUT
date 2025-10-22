package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"matcher/internal/audit"
	"matcher/internal/config"
	"matcher/internal/db"
	"matcher/internal/matcher"
	"matcher/internal/repo"
	"matcher/internal/router"
	"matcher/internal/routing"
	"matcher/internal/seed"
	"matcher/internal/worker"
)

func main() {
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	pool, err := db.New(ctx, cfg.DBURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pool.Close()

	if err := applyMigrations(ctx, pool); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	if cfg.SeedPoints {
		if err := seed.SeedLogisticPoints(ctx, pool, "data/logistic_points.json"); err != nil {
			log.Printf("seed logistic_points: %v", err)
		} else {
			log.Printf("seeded logistic_points")
		}
	}

	repository := repo.New(pool)

	var routerClient routing.RouterClient
	if cfg.UseYandex && cfg.YandexAPIKey != "" {
		routerClient = routing.YandexRouter{
			BaseURL:  cfg.YandexMatrixURL,
			APIKey:   cfg.YandexAPIKey,
			MaxBatch: 25,
		}
		log.Printf("routing provider: yandex")
	} else {
		routerClient = routing.HaversineRouter{AvgSpeedKmh: cfg.MatchAvgSpeedKmh}
		log.Printf("routing provider: haversine fallback")
	}

	var auditClient *audit.Client
	if cfg.AuditURL != "" {
		auditClient = &audit.Client{BaseURL: cfg.AuditURL, From: "matcher"}
	}

	m := matcher.NewService(repository, routerClient, auditClient, cfg.MatchMaxDetourMin, cfg.MatchRadiusKm, cfg.MatchAvgSpeedKmh)

	// HTTP
	handler := router.New(repository, m)
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// worker — запускаем только если интервал > 0 (можно выключить на время ручных тестов)
	if cfg.WorkerInterval > 0 {
		w := &worker.Worker{Repo: repository, Matcher: m, Interval: cfg.WorkerInterval}
		go w.Run(ctx)
	} else {
		log.Printf("worker disabled")
	}

	go func() {
		log.Printf("be-1 (matcher) listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")
	shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shCtx)
}

func applyMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	b, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		return err
	}
	for _, stmt := range strings.Split(string(b), ";\n") {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := pool.Exec(ctx, stmt); err != nil {
			return err
		}
	}
	return nil
}
