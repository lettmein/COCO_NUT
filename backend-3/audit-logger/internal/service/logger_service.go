package service

import (
	"context"
	"myapp/internal/domain"
)

type LoggerService struct {
	repo domain.LoggerRepository
}

func NewLoggerService(repo domain.LoggerRepository) domain.LoggerService {
	return &LoggerService{repo: repo}
}

func (s *LoggerService) Create(ctx context.Context, log *domain.Log) error {

	return s.repo.Create(ctx, log)
}

func (s *LoggerService) GetAllWithFilter(ctx context.Context, filter domain.LogFilter) ([]domain.Log, error) {
	logs, err := s.repo.GetWithFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Для пагинации можно добавить получение общего количества
	// total, err := s.repo.GetCountWithFilter(ctx, filter)
	// if err != nil {
	//     return nil, err
	// }

	// Пока возвращаем без общего количества или можно добавить позже
	return logs, nil
}
