package server

import (
	"database/sql"
	"fmt"
	"git.a7ru.app/a7hack/coco-nut/backend-2/office/internal/client"
	"git.a7ru.app/a7hack/coco-nut/backend-2/office/internal/config"
	"git.a7ru.app/a7hack/coco-nut/backend-2/office/internal/controller"
	"git.a7ru.app/a7hack/coco-nut/backend-2/office/internal/repo"
	"git.a7ru.app/a7hack/coco-nut/backend-2/office/internal/service"
	"git.a7ru.app/a7hack/coco-nut/backend-2/office/internal/utils"
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
		AppName: "Office Service",
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

	routeRepo := repo.NewRouteRepository(s.db)
	routeService := service.NewRouteService(routeRepo)
	routeController := controller.NewRouteController(routeService)

	api := s.app.Group("/api/v1")

	routes := api.Group("/routes")
	routes.Post("/", routeController.CreateRoute)
	routes.Get("/", routeController.GetAllRoutes)
	routes.Get("/:id", routeController.GetRoute)
	routes.Patch("/:id/status", routeController.UpdateRouteStatus)
	routes.Delete("/:id", routeController.DeleteRoute)

	s.app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "office-service",
		})
	})
}

func (s *Server) Start() error {
	s.SetupRoutes()
	addr := fmt.Sprintf(":%s", s.config.Port)
	log.Printf("Starting Office Service on port %s", s.config.Port)
	return s.app.Listen(addr)
}
