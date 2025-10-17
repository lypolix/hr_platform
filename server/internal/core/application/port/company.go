package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/hr-platform-mosprom/internal/core/domain"
)




// Репозиторий компаний
type CompanyRepository interface {
	Save(ctx context.Context, c *domain.Company) error
	ByID(ctx context.Context, id uuid.UUID) (*domain.Company, error)
	ByINN(ctx context.Context, inn string) (*domain.Company, error)
	ByLogin(ctx context.Context, login string) (*domain.Company, error)
	Approved(ctx context.Context, limit, offset int) ([]*domain.Company, error)
	ByRepresentative(ctx context.Context, repID uuid.UUID) ([]*domain.Company, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// Сервис компаний
type CompanyService interface {
	Register(ctx context.Context, in RegisterCompanyInput) (*domain.Company, error)
	Approve(ctx context.Context, id uuid.UUID) error
	Login(ctx context.Context, login, password string) (*domain.Company, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, in UpdateCompanyInput) (*domain.Company, error)
	ChangeCredentials(ctx context.Context, id uuid.UUID, login, password string) (*domain.Company, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Company, error)
	ListApproved(ctx context.Context, limit, offset int) ([]*domain.Company, error)
	ListByRepresentative(ctx context.Context, repID uuid.UUID) ([]*domain.Company, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type RegisterCompanyInput struct {
	Title            string
	Description      string
	Contacts         string
	INN              string
	Address          string
	Website          string
	LogoURL          string
	RepresentativeID uuid.UUID
	Login            string
	Password         string
}

type UpdateCompanyInput struct {
	Title       string
	Description string
	Contacts    string
	Address     string
	Website     string
	LogoURL     string
}
