package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/hr-platform-mosprom/internal/core/domain"
)

type ResponseRepository interface {
	Create(ctx context.Context, resp *domain.Response) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Response, error)
	GetByVacancyID(ctx context.Context, vacancyID uuid.UUID, limit, offset int) ([]*domain.Response, error)
	GetWithFilters(ctx context.Context, filter domain.ResponseFilter) ([]*domain.Response, error)
	Update(ctx context.Context, resp *domain.Response) error
	Delete(ctx context.Context, id uuid.UUID) error
}


type ResponseService interface {
	Create(ctx context.Context, req domain.CreateResponseRequest) (*domain.Response, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Response, error)
	ListByVacancy(ctx context.Context, vacancyID uuid.UUID, limit, offset int) ([]*domain.Response, error)
	Search(ctx context.Context, filter domain.ResponseFilter) ([]*domain.Response, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	Delete(ctx context.Context, id uuid.UUID) error
}
