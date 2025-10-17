package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/hr-platform-mosprom/internal/core/domain"
)

type VacancyRepository interface {
	Create(ctx context.Context, vacancy *domain.Vacancy) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Vacancy, error)
	GetByCompanyID(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*domain.Vacancy, error)
	GetAll(ctx context.Context, limit, offset int) ([]*domain.Vacancy, error)
	GetWithFilters(ctx context.Context, filter domain.VacancyFilter) ([]*domain.Vacancy, error)
	Update(ctx context.Context, vacancy *domain.Vacancy) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type VacancyService interface {
	Create(ctx context.Context, req domain.CreateVacancyRequest) (*domain.Vacancy, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Vacancy, error)
	ListByCompany(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*domain.Vacancy, error)
	ListAll(ctx context.Context, limit, offset int) ([]*domain.Vacancy, error)
	Search(ctx context.Context, filter domain.VacancyFilter) ([]*domain.Vacancy, error)
	Update(ctx context.Context, id uuid.UUID, req domain.UpdateVacancyRequest) (*domain.Vacancy, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Activate(ctx context.Context, id uuid.UUID) error
	Deactivate(ctx context.Context, id uuid.UUID) error
}
