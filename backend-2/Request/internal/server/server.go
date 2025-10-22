package server

import (
	"database/sql"
	"fmt"
	"git.a7ru.app/a7hack/coco-nut/backend-2/request/internal/client"
	"git.a7ru.app/a7hack/coco-nut/backend-2/request/internal/config"
	"git.a7ru.app/a7hack/coco-nut/backend-2/request/internal/controller"
	"git.a7ru.app/a7hack/coco-nut/backend-2/request/internal/repo"
	"git.a7ru.app/a7hack/coco-nut/backend-2/request/internal/service"
	"git.a7ru.app/a7hack/coco-nut/backend-2/request/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
)

type Server struct {
	app    *fiber.App
	config *config.Config
	db     *sql.DB
}

func NewServer(cfg *config.Config, db *sql.DB) *Server {
	app := fiber.New(fiber.Config{
		AppName: "Request Service",
	})

	return &Server{
		app:    app,
		config: cfg,
		db:     db,
	}
}

func (s *Server) SetupRoutes() {
	auditClient := client.NewAuditClient(s.config.AuditServiceURL)

	s.app.Use(cors.New())
	s.app.Use(logger.New())
	s.app.Use(utils.AuditMiddleware(s.config.ServiceName, auditClient))

	requestRepo := repo.NewRequestRepository(s.db)
	requestService := service.NewRequestService(requestRepo)
	requestController := controller.NewRequestController(requestService)

	api := s.app.Group("/api/v1")

	requests := api.Group("/requests")
	requests.Post("/", requestController.CreateRequest)
	requests.Get("/", requestController.GetAllRequests)
	requests.Get("/:id", requestController.GetRequest)
	requests.Patch("/:id/status", requestController.UpdateRequestStatus)
	requests.Delete("/:id", requestController.DeleteRequest)

	s.app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"service": "request-service",
		})
	})
}

func (s *Server) Start() error {
	s.SetupRoutes()
	addr := fmt.Sprintf(":%s", s.config.Port)
	log.Printf("Starting Request Service on port %s", s.config.Port)
	return s.app.Listen(addr)
}
