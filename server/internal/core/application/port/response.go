package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/hr-platform-mosprom/internal/core/domain"
)

type ResponseRepository interface {
	Save(ctx context.Context, r *domain.Response) error
	ByID(ctx context.Context, id uuid.UUID) (*domain.Response, error)
	ByVacancy(ctx context.Context, vacancyID uuid.UUID, limit, offset int) ([]*domain.Response, error)
	Search(ctx context.Context, f domain.ResponseFilter) ([]*domain.Response, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type ResponseService interface {
	Create(ctx context.Context, in CreateResponseInput) (*domain.Response, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Response, error)
	ListByVacancy(ctx context.Context, vacancyID uuid.UUID, limit, offset int) ([]*domain.Response, error)
	Search(ctx context.Context, f domain.ResponseFilter) ([]*domain.Response, error)
	SetStatus(ctx context.Context, id uuid.UUID, status string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateResponseInput struct {
	VacancyID   uuid.UUID
	FullName    string
	Email       string
	Phone       string
	CoverLetter string
	ResumeURL   string
}
