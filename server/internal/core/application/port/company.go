package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/hr-platform-mosprom/internal/core/domain"
)

type CompanyRepository interface {
	Create(ctx context.Context, company *domain.Company) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Company, error)
	GetByINN(ctx context.Context, inn string) (*domain.Company, error)
	GetByLogin(ctx context.Context, login string) (*domain.Company, error)
	GetByRepresentativeID(ctx context.Context, repID uuid.UUID) ([]*domain.Company, error)
	GetAllApproved(ctx context.Context, limit, offset int) ([]*domain.Company, error)
	GetPending(ctx context.Context, limit, offset int) ([]*domain.Company, error)
	Update(ctx context.Context, company *domain.Company) error
	Delete(ctx context.Context, id uuid.UUID) error
	SetApproved(ctx context.Context, id uuid.UUID, approved bool) error
}


type CompanyService interface {
	Register(ctx context.Context, req domain.CreateCompanyRequest) (*domain.Company, error)
	Approve(ctx context.Context, id uuid.UUID) error
	Login(ctx context.Context, login, password string) (*domain.Company, error)

	GetByID(ctx context.Context, id uuid.UUID) (*domain.Company, error)
	ListApproved(ctx context.Context, limit, offset int) ([]*domain.Company, error)
	ListByRepresentative(ctx context.Context, repID uuid.UUID) ([]*domain.Company, error)
	Update(ctx context.Context, id uuid.UUID, req domain.UpdateCompanyRequest) (*domain.Company, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
