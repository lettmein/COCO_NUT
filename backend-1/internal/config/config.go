package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DBURL             string
	UseYandex         bool
	YandexAPIKey      string
	YandexMatrixURL   string
	MatchRadiusKm     float64
	MatchMaxDetourMin float64
	MatchAvgSpeedKmh  float64
	WorkerInterval    time.Duration
	AuditURL          string
	SeedPoints        bool
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
func getBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		if v == "1" || v == "true" || v == "TRUE" || v == "yes" {
			return true
		}
		if v == "0" || v == "false" || v == "FALSE" || v == "no" {
			return false
		}
	}
	return def
}
func getFloat(key string, def float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return def
}
func getInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

func Load() Config {
	return Config{
		DBURL:             getEnv("DATABASE_URL", "postgres://postgres:postgres@db:5432/app?sslmode=disable"),
		UseYandex:         getBool("USE_YANDEX", false),
		YandexAPIKey:      getEnv("YANDEX_API_KEY", ""),
		YandexMatrixURL:   getEnv("YANDEX_MATRIX_URL", "https://api.routing.yandex.net/v2/matrix"),
		MatchRadiusKm:     getFloat("MATCH_RADIUS_KM", 30),
		MatchMaxDetourMin: getFloat("MATCH_MAX_DETOUR_MIN", 60),
		MatchAvgSpeedKmh:  getFloat("MATCH_AVG_SPEED_KMH", 60),
		WorkerInterval:    time.Duration(getInt("MATCH_WORKER_INTERVAL_SEC", 30)) * time.Second,
		AuditURL:          getEnv("AUDIT_URL", ""),
		SeedPoints:        getBool("SEED_POINTS", false),
	}
}
