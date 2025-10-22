package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"matcher/internal/matcher"
	"matcher/internal/repo"
)

type Server struct {
	R *repo.Repo
	M *matcher.Service
}

func New(r *repo.Repo, m *matcher.Service) http.Handler {
	s := &Server{R: r, M: m}
	rt := chi.NewRouter()

	rt.Get("/health", s.health)
	rt.Post("/routes/{id}/match", s.matchRoute)
	rt.Get("/routes/{id}/assignments", s.getAssignments)

	return rt
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) matchRoute(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	routeID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "bad route id", http.StatusBadRequest)
		return
	}
	as, err := s.M.MatchRoute(r.Context(), routeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"route_id": routeID, "assigned": as})
}

func (s *Server) getAssignments(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	routeID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "bad route id", http.StatusBadRequest)
		return
	}
	as, err := s.R.GetAssignments(r.Context(), routeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"route_id": routeID, "assignments": as})
}
