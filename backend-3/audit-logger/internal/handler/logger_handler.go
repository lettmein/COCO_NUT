package handler

import (
	"encoding/json"
	"myapp/internal/domain"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type LoggerHandler struct {
	service domain.LoggerService
}

func NewLoggerHandler(service domain.LoggerService) *LoggerHandler {
	return &LoggerHandler{service: service}
}

func (h *LoggerHandler) Start(address string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", h.healthCheck)
	mux.HandleFunc("/users/{id}", h.getUser)

	h.server = &http.Server{
		Addr:    address,
		Handler: mux,
	}

	h.logger.Infof("Starting HTTP server on %s", address)
	if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		h.logger.Fatalf("HTTP server error: %v", err)
	}
}

// CreateLog обрабатывает создание новой записи лога
// @Summary Создать запись лога
// @Description Создает новую запись в логах
// @Tags logs
// @Accept json
// @Produce json
// @Param log body domain.Log true "Данные лога"
// @Success 201 {object} domain.Log
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /logs [post]
func (h *LoggerHandler) CreateLog(w http.ResponseWriter, r *http.Request) {
	var log domain.Log
	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Устанавливаем время создания если не указано
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}

	if err := h.service.Create(r.Context(), &log); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create log: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, log)
}

// GetLogs обрабатывает получение логов с фильтрацией
// @Summary Получить логи с фильтрацией
// @Description Возвращает список логов с возможностью фильтрации
// @Tags logs
// @Produce json
// @Param service_type query string false "Тип сервиса"
// @Param uri query string false "URI"
// @Param start_time query string false "Начальное время (RFC3339)"
// @Param end_time query string false "Конечное время (RFC3339)"
// @Param min_duration query number false "Минимальная длительность"
// @Param max_duration query number false "Максимальная длительность"
// @Param limit query int false "Лимит (по умолчанию 50, максимум 1000)"
// @Param offset query int false "Смещение (по умолчанию 0)"
// @Success 200 {array} domain.Log
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /logs [get]
func (h *LoggerHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	filter := domain.LogFilter{}

	// Парсим query параметры
	query := r.URL.Query()

	if serviceType := query.Get("service_type"); serviceType != "" {
		filter.ServiceType = serviceType
	}

	if uri := query.Get("uri"); uri != "" {
		filter.URI = uri
	}

	if startTimeStr := query.Get("start_time"); startTimeStr != "" {
		if startTime, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			filter.StartTime = startTime
		} else {
			respondWithError(w, http.StatusBadRequest, "Invalid start_time format, use RFC3339")
			return
		}
	}

	if endTimeStr := query.Get("end_time"); endTimeStr != "" {
		if endTime, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			filter.EndTime = endTime
		} else {
			respondWithError(w, http.StatusBadRequest, "Invalid end_time format, use RFC3339")
			return
		}
	}

	if minDurationStr := query.Get("min_duration"); minDurationStr != "" {
		if minDuration, err := strconv.ParseFloat(minDurationStr, 64); err == nil {
			filter.MinDuration = minDuration
		} else {
			respondWithError(w, http.StatusBadRequest, "Invalid min_duration format")
			return
		}
	}

	if maxDurationStr := query.Get("max_duration"); maxDurationStr != "" {
		if maxDuration, err := strconv.ParseFloat(maxDurationStr, 64); err == nil {
			filter.MaxDuration = maxDuration
		} else {
			respondWithError(w, http.StatusBadRequest, "Invalid max_duration format")
			return
		}
	}

	if limitStr := query.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		} else {
			respondWithError(w, http.StatusBadRequest, "Invalid limit format")
			return
		}
	} else {
		filter.Limit = 50 // значение по умолчанию
	}

	if offsetStr := query.Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		} else {
			respondWithError(w, http.StatusBadRequest, "Invalid offset format")
			return
		}
	}

	// Ограничиваем максимальный лимит
	if filter.Limit > 1000 {
		filter.Limit = 1000
	}

	logs, err := h.service.GetAllWithFilter(r.Context(), filter)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get logs: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, logs)
}

// GetSlowLogs обрабатывает получение медленных запросов
// @Summary Получить медленные запросы
// @Description Возвращает список логов с длительностью выше указанной
// @Tags logs
// @Produce json
// @Param min_duration query number false "Минимальная длительность (по умолчанию 1.0)"
// @Param limit query int false "Лимит (по умолчанию 100)"
// @Success 200 {array} domain.Log
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /logs/slow [get]
func (h *LoggerHandler) GetSlowLogs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	minDuration := 1.0 // значение по умолчанию
	if minDurationStr := query.Get("min_duration"); minDurationStr != "" {
		if parsed, err := strconv.ParseFloat(minDurationStr, 64); err == nil {
			minDuration = parsed
		} else {
			respondWithError(w, http.StatusBadRequest, "Invalid min_duration format")
			return
		}
	}

	limit := 100 // значение по умолчанию
	if limitStr := query.Get("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		} else {
			respondWithError(w, http.StatusBadRequest, "Invalid limit format")
			return
		}
	}

	logs, err := h.service.GetSlowRequests(r.Context(), minDuration, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get slow logs: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, logs)
}

// GetLogsByService обрабатывает получение логов по типу сервиса
// @Summary Получить логи по типу сервиса
// @Description Возвращает список логов для указанного типа сервиса
// @Tags logs
// @Produce json
// @Param service_type path string true "Тип сервиса"
// @Param limit query int false "Лимит (по умолчанию 100)"
// @Success 200 {array} domain.Log
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /logs/service/{service_type} [get]
func (h *LoggerHandler) GetLogsByService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceType := vars["service_type"]

	if serviceType == "" {
		respondWithError(w, http.StatusBadRequest, "Service type is required")
		return
	}

	query := r.URL.Query()
	limit := 100 // значение по умолчанию
	if limitStr := query.Get("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		} else {
			respondWithError(w, http.StatusBadRequest, "Invalid limit format")
			return
		}
	}

	logs, err := h.service.GetRequestsByService(r.Context(), serviceType, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get service logs: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, logs)
}

// HealthCheck обрабатывает проверку здоровья сервиса
// @Summary Проверка здоровья
// @Description Проверяет доступность сервиса
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *LoggerHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok", "timestamp": time.Now().Format(time.RFC3339)})
}

// Вспомогательные функции для ответов

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ErrorResponse{
		Error:   http.StatusText(code),
		Message: message,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal Server Error", "message": "Failed to marshal response"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
