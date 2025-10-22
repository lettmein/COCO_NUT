package handler

import (
	"myapp/internal/domain"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRouter(service domain.LoggerService) *mux.Router {
	handler := NewLoggerHandler(service)

	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Log routes
	api.HandleFunc("/logs", handler.CreateLog).Methods("POST")
	api.HandleFunc("/logs", handler.GetLogs).Methods("GET")

	return router
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Логирование запросов
		next.ServeHTTP(w, r)
	})
}

func contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
