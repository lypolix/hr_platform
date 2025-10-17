package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/hr-platform-mosprom/internal/core/domain"
)

type VacancyRepository interface {
	Save(ctx context.Context, v *domain.Vacancy) error
	ByID(ctx context.Context, id uuid.UUID) (*domain.Vacancy, error)
	ByCompany(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*domain.Vacancy, error)
	All(ctx context.Context, limit, offset int) ([]*domain.Vacancy, error)
	Search(ctx context.Context, f domain.VacancyFilter) ([]*domain.Vacancy, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type VacancyService interface {
	Create(ctx context.Context, in CreateVacancyInput) (*domain.Vacancy, error)
	Update(ctx context.Context, id uuid.UUID, in UpdateVacancyInput) (*domain.Vacancy, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Vacancy, error)
	ListByCompany(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*domain.Vacancy, error)
	ListAll(ctx context.Context, limit, offset int) ([]*domain.Vacancy, error)
	Search(ctx context.Context, f domain.VacancyFilter) ([]*domain.Vacancy, error)
	Activate(ctx context.Context, id uuid.UUID) error
	Deactivate(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateVacancyInput struct {
	CompanyID       uuid.UUID
	Title           string
	Description     string
	Contacts        string
	Requirements     string
	Responsibilities string
	Conditions       string
	SalaryFrom       *int
	SalaryTo         *int
	Employment       string
	Schedule         string
	Experience       string
	Education        string
	Location         string
}

type UpdateVacancyInput struct {
	Title           string
	Description     string
	Contacts        string
	Requirements     string
	Responsibilities string
	Conditions       string
	SalaryFrom       *int
	SalaryTo         *int
	Employment       string
	Schedule         string
	Experience       string
	Education        string
	Location         string
	IsActive         *bool
}
