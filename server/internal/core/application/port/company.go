package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/hr-platform-mosprom/internal/core/domain"
)

// Репозиторий компаний
type CompanyRepository interface {
	Save(ctx context.Context, company *domain.Company) error
	GetByLogin(ctx context.Context, login string) (*domain.Company, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Company, error)
	GetByINN(ctx context.Context, inn string) (*domain.Company, error)
}

// Входные данные для регистрации компании
type SignUpCompanyData struct {
	Login            string
	Password         string
	Title            string
	INN              string
	Description      string
	Contacts         string
	Address          string
	RepresentativeID uuid.UUID
}

// Входные данные для авторизации
type SignInCompanyData struct {
	Login    string
	Password string
}

// Результат для компании (как UniversityResult)
type CompanyResult struct {
	Login string
	Title string
	INN   string
}

// Результат с токеном (как UniversityWithTokenResult)
type CompanyWithTokenResult struct {
	CompanyResult
	Token string
}

// Данные для обновления профиля компании
type UpdateCompanyProfileData struct {
	Title       string
	Description string
	Contacts    string
	Address     string
}

// Данные для смены учётных данных компании
type ChangeCompanyCredentialsData struct {
	Login    string
	Password string
}
