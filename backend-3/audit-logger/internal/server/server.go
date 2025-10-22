package server

import (
	"context"
	"net/http"
	"time"

	"myapp/internal/config"
	"myapp/internal/router"
	"myapp/pkg/logger"
)

type HTTPServer struct {
	server *http.Server
	config *config.Config
	log    *logger.Logger
}

func NewHTTPServer(router *router.Router, config *config.Config, log *logger.Logger) *HTTPServer {
	return &HTTPServer{
		server: &http.Server{
			Addr:         config.HTTP.Address,
			Handler:      router.GetHandler(),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		config: config,
		log:    log,
	}
}

func (s *HTTPServer) Start() {
	s.log.Infof("Starting HTTP server on %s", s.config.HTTP.Address)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.log.Fatalf("Failed to start server: %v", err)
	}
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	s.log.Info("Shutting down HTTP server")
	return s.server.Shutdown(ctx)
}
