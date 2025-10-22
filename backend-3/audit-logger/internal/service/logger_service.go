package service

import (
	"context"
	"myapp/internal/domain"

	"time"
)

type loggerService struct {
	repo domain.LoggerRepository
}

func NewLoggerService(repo domain.LoggerRepository) domain.LoggerService {
	return &loggerService{repo: repo}
}

func (s *loggerService) Create(ctx context.Context, log *domain.Log) error {
	// Валидация обязательных полей
	if log.ReqServiceType == "" {
		return domain.ErrInvalidServiceType
	}
	if log.Uri == "" {
		return domain.ErrInvalidURI
	}
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}

	return s.repo.Create(ctx, log)
}

func (s *loggerService) GetAllWithFilter(ctx context.Context, filter domain.LogFilter) ([]*domain.Log, error) {
	// Применяем значения по умолчанию для пагинации
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	if filter.Limit > 1000 {
		filter.Limit = 1000
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	// Если фильтр пустой, используем базовый метод GetAll
	if s.isFilterEmpty(filter) {
		return s.repo.GetAll(ctx, filter.Limit, filter.Offset)
	}

	// Используем метод с фильтрацией (если он реализован в репозитории)
	if repoWithFilter, ok := s.repo.(interface {
		GetWithFilter(ctx context.Context, filter domain.LogFilter) ([]*domain.Log, error)
	}); ok {
		return repoWithFilter.GetWithFilter(ctx, filter)
	}

	// Альтернативная реализация - получаем все и фильтруем в памяти
	// (не рекомендуется для больших объемов данных)
	return s.getAllAndFilterInMemory(ctx, filter)
}

func (s *loggerService) isFilterEmpty(filter domain.LogFilter) bool {
	return filter.ServiceType == "" &&
		filter.URI == "" &&
		filter.StartTime.IsZero() &&
		filter.EndTime.IsZero() &&
		filter.MinDuration <= 0 &&
		filter.MaxDuration <= 0
}

func (s *loggerService) getAllAndFilterInMemory(ctx context.Context, filter domain.LogFilter) ([]*domain.Log, error) {
	// Получаем все записи (или большую партию)
	logs, err := s.repo.GetAll(ctx, 10000, 0) // Ограничиваем для безопасности
	if err != nil {
		return nil, err
	}

	var filteredLogs []*domain.Log
	for _, log := range logs {
		if s.matchesFilter(log, filter) {
			filteredLogs = append(filteredLogs, log)
		}
	}

	// Применяем пагинацию
	start := filter.Offset
	if start > len(filteredLogs) {
		start = len(filteredLogs)
	}
	end := start + filter.Limit
	if end > len(filteredLogs) {
		end = len(filteredLogs)
	}

	return filteredLogs[start:end], nil
}

func (s *loggerService) matchesFilter(log *domain.Log, filter domain.LogFilter) bool {
	// Фильтр по типу сервиса
	if filter.ServiceType != "" {
		if log.ReqServiceType != filter.ServiceType && log.RespServiceType != filter.ServiceType {
			return false
		}
	}

	// Фильтр по URI
	if filter.URI != "" {
		if !containsSubstring(log.Uri, filter.URI) {
			return false
		}
	}

	// Фильтр по времени начала
	if !filter.StartTime.IsZero() {
		if log.CreatedAt.Before(filter.StartTime) {
			return false
		}
	}

	// Фильтр по времени окончания
	if !filter.EndTime.IsZero() {
		if log.CreatedAt.After(filter.EndTime) {
			return false
		}
	}

	// Фильтр по минимальному времени выполнения
	if filter.MinDuration > 0 {
		if log.DurationTime < filter.MinDuration {
			return false
		}
	}

	// Фильтр по максимальному времени выполнения
	if filter.MaxDuration > 0 {
		if log.DurationTime > filter.MaxDuration {
			return false
		}
	}

	return true
}

func containsSubstring(str, substr string) bool {
	if len(substr) > len(str) {
		return false
	}
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Дополнительные методы сервиса
func (s *loggerService) GetByID(ctx context.Context, id int) (*domain.Log, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidID
	}
	return s.repo.GetByID(ctx, id)
}

func (s *loggerService) GetSlowRequests(ctx context.Context, minDuration float64, limit int) ([]*domain.Log, error) {
	if minDuration <= 0 {
		minDuration = 1.0 // секунда по умолчанию
	}
	if limit <= 0 {
		limit = 100
	}

	filter := domain.LogFilter{
		MinDuration: minDuration,
		Limit:       limit,
	}

	return s.GetAllWithFilter(ctx, filter)
}

func (s *loggerService) GetRequestsByService(ctx context.Context, serviceType string, limit int) ([]*domain.Log, error) {
	if serviceType == "" {
		return nil, domain.ErrInvalidServiceType
	}
	if limit <= 0 {
		limit = 100
	}

	filter := domain.LogFilter{
		ServiceType: serviceType,
		Limit:       limit,
	}

	return s.GetAllWithFilter(ctx, filter)
}

func (s *loggerService) GetRequestsByTimeRange(ctx context.Context, startTime, endTime time.Time, limit int) ([]*domain.Log, error) {
	if startTime.After(endTime) {
		return nil, domain.ErrInvalidTimeRange
	}
	if limit <= 0 {
		limit = 100
	}

	filter := domain.LogFilter{
		StartTime: startTime,
		EndTime:   endTime,
		Limit:     limit,
	}

	return s.GetAllWithFilter(ctx, filter)
}
