package service

import (
	"git.a7ru.app/a7hack/coco-nut/backend-2/request/internal/dto"
	"git.a7ru.app/a7hack/coco-nut/backend-2/request/internal/repo"
)

type RequestService struct {
	repo *repo.RequestRepository
}

func NewRequestService(repo *repo.RequestRepository) *RequestService {
	return &RequestService{repo: repo}
}

func (s *RequestService) CreateRequest(req *dto.CreateRequestDTO) (*dto.RequestResponse, error) {
	return s.repo.Create(req)
}

func (s *RequestService) GetRequest(id string) (*dto.RequestResponse, error) {
	return s.repo.GetByID(id)
}

func (s *RequestService) GetAllRequests() ([]*dto.RequestResponse, error) {
	return s.repo.GetAll()
}

func (s *RequestService) UpdateRequestStatus(id string, status string) error {
	return s.repo.UpdateStatus(id, status)
}

func (s *RequestService) DeleteRequest(id string) error {
	return s.repo.Delete(id)
}

func (s *RequestService) GetRequestsByStatus(status string) ([]*dto.RequestResponse, error) {
	return s.repo.GetByStatus(status)
}
