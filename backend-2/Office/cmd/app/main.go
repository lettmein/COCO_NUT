package main

import (
	"database/sql"
	"git.a7ru.app/a7hack/coco-nut/backend-2/office/internal/config"
	"git.a7ru.app/a7hack/coco-nut/backend-2/office/internal/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	godotenv.Load()

	cfg := config.NewConfig()

	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to database")

	srv := server.NewServer(cfg, db)

	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
