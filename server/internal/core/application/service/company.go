package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hr-platform-mosprom/internal/core/application/port"
	"github.com/hr-platform-mosprom/internal/core/domain"
	"golang.org/x/crypto/bcrypt"
)

type companyService struct {
	repo port.CompanyRepository
}

func NewCompanyService(repo port.CompanyRepository) port.CompanyService {
	return &companyService{repo: repo}
}


func (s *companyService) Register(ctx context.Context, req domain.CreateCompanyRequest) (*domain.Company, error) {
	if existing, _ := s.repo.GetByINN(ctx, req.INN); existing != nil {
		return nil, fmt.Errorf("company with INN %s already exists", req.INN)
	}
	if existedLogin, _ := s.repo.GetByLogin(ctx, req.Login); existedLogin != nil {
		return nil, fmt.Errorf("login %s already taken", req.Login)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	cmp := domain.NewCompany(req, string(hash))
	if err := s.repo.Create(ctx, cmp); err != nil {
		return nil, fmt.Errorf("create company: %w", err)
	}
	return cmp, nil
}

func (s *companyService) Login(ctx context.Context, login, password string) (*domain.Company, error) {
	cmp, err := s.repo.GetByLogin(ctx, login)
	if err != nil || cmp == nil {
		return nil, errors.New("invalid credentials")
	}
	if bcrypt.CompareHashAndPassword([]byte(cmp.PasswordHash), []byte(password)) != nil {
		return nil, errors.New("invalid credentials")
	}
	if !cmp.Approved {
		return nil, errors.New("company is not approved yet")
	}
	return cmp, nil
}

func (s *companyService) Approve(ctx context.Context, id uuid.UUID) error {
	return s.repo.SetApproved(ctx, id, true)
}

func (s *companyService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *companyService) ListApproved(ctx context.Context, limit, offset int) ([]*domain.Company, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.GetAllApproved(ctx, limit, offset)
}

func (s *companyService) ListByRepresentative(ctx context.Context, repID uuid.UUID) ([]*domain.Company, error) {
	return s.repo.GetByRepresentativeID(ctx, repID)
}

func (s *companyService) Update(ctx context.Context, id uuid.UUID, req domain.UpdateCompanyRequest) (*domain.Company, error) {
	cmp, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get company: %w", err)
	}
	cmp.Update(req)
	if err := s.repo.Update(ctx, cmp); err != nil {
		return nil, fmt.Errorf("update company: %w", err)
	}
	return cmp, nil
}

func (s *companyService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
