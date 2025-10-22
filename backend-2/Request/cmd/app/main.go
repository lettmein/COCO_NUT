package main

import (
	"database/sql"
	"git.a7ru.app/a7hack/coco-nut/backend-2/request/internal/config"
	"git.a7ru.app/a7hack/coco-nut/backend-2/request/internal/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
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

	log.Println("Running migrations...")
	if err := goose.Up(db, "/root/migration"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations completed successfully")

	srv := server.NewServer(cfg, db)

	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
