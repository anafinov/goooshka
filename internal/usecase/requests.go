package usecase

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

type RequestsUseCase struct {
	repo      RequestRepository
	publisher MessagePublisher
}

func NewRequestsUseCase(repo RequestRepository, publisher MessagePublisher) *RequestsUseCase {
	return &RequestsUseCase{
		repo:      repo,
		publisher: publisher,
	}
}

func (u *RequestsUseCase) CreateRequest(userID int, reqType string) (*domain.Request, error) {
	req, err := u.repo.CreateRequest(userID, reqType)
	if err != nil {
		return nil, err
	}

	// Publish to RabbitMQ
	if err := u.publisher.PublishRequest(req); err != nil {
		// Log error, but maybe don't fail the HTTP request, or implement a retry/outbox pattern
		// For simplicity, we just return the error here
		return nil, err
	}

	return req, nil
}

func (u *RequestsUseCase) GetMyRequests(userID int) ([]domain.Request, error) {
	return u.repo.GetRequestsByUserID(userID)
}
