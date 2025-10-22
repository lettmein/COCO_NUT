package router

import (
	"myapp/internal/handler"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	router        *mux.Router
	loggerHandler *handler.LoggerHandler
}

func NewRouter(loggerHandler *handler.LoggerHandler) *Router {
	r := &Router{
		router:        mux.NewRouter(),
		loggerHandler: loggerHandler,
	}

	r.setupRoutes()
	return r
}

func (r *Router) setupRoutes() {
	// API v1 routes
	api := r.router.PathPrefix("/api/v1").Subrouter()

	// Logs routes
	logs := api.PathPrefix("/logs").Subrouter()
	logs.HandleFunc("", r.loggerHandler.CreateLog).Methods("POST")
	logs.HandleFunc("", r.loggerHandler.GetLogs).Methods("GET")

}

func (r *Router) GetHandler() http.Handler {
	return r.router
}
