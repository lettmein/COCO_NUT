package handler

import (
	"encoding/json"
	"myapp/internal/domain"
	"myapp/pkg/logger"
	"net/http"
)

type LoggerHandler struct {
	service domain.LoggerService
	log     *logger.Logger
}

func NewLoggerHandler(service domain.LoggerService, log *logger.Logger) *LoggerHandler {
	return &LoggerHandler{
		service: service,
		log:     log,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

func (h *LoggerHandler) CreateLog(w http.ResponseWriter, r *http.Request) {
	var request domain.Log

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.service.Create(r.Context(), &request)

	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusCreated, "log created successfully")
}

func (h *LoggerHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := h.service.GetAllWithFilter(r.Context(), domain.LogFilter{})
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, logs)
}

// Вспомогательные методы
func (h *LoggerHandler) respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.log.Errorf("Failed to encode JSON response: %v", err)
	}
}

func (h *LoggerHandler) respondWithError(w http.ResponseWriter, statusCode int, message string) {
	h.respondWithJSON(w, statusCode, ErrorResponse{Error: message})
}
