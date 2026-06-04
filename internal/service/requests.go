package service

import (
	"student-cert-service/internal/domain"
)

type RequestRepository interface {
	CreateRequest(userID int, reqType string) (*domain.Request, error)
	GetRequestsByUserID(userID int) ([]domain.Request, error)
}

type MessagePublisher interface {
	PublishRequest(req *domain.Request) error
}

type RequestsService struct {
	repo      RequestRepository
	publisher MessagePublisher
}

func NewRequestsService(repo RequestRepository, publisher MessagePublisher) *RequestsService {
	return &RequestsService{
		repo:      repo,
		publisher: publisher,
	}
}

func (u *RequestsService) CreateRequest(userID int, reqType string) (*domain.Request, error) {
	req, err := u.repo.CreateRequest(userID, reqType)
	if err != nil {
		return nil, err
	}

	if err := u.publisher.PublishRequest(req); err != nil {

		return nil, err
	}

	return req, nil
}

func (u *RequestsService) GetMyRequests(userID int) ([]domain.Request, error) {
	return u.repo.GetRequestsByUserID(userID)
}
