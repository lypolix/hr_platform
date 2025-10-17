package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hr-platform-mosprom/internal/core/application/port"
	"github.com/hr-platform-mosprom/internal/core/domain"
)

type companyService struct {
	repo  port.CompanyRepository
	pass  port.PasswordService
	clock port.Clock
}

func NewCompanyService(r port.CompanyRepository, p port.PasswordService, c port.Clock) *companyService {
	return &companyService{repo: r, pass: p, clock: c}
}

func (s *companyService) Register(ctx context.Context, in port.RegisterCompanyInput) (*domain.Company, error) {
	if existing, _ := s.repo.ByINN(ctx, in.INN); existing != nil {
		return nil, fmt.Errorf("company with INN exists")
	}
	if existing, _ := s.repo.ByLogin(ctx, in.Login); existing != nil {
		return nil, fmt.Errorf("login already taken")
	}
	hash, err := s.pass.Hash(in.Password)
	if err != nil {
		return nil, err
	}
	now := s.clock.Now()
	c, err := domain.CreateCompany(domain.CreateCompanyAttrs{
		Title:            in.Title,
		Description:      in.Description,
		Contacts:         in.Contacts,
		INN:              in.INN,
		Address:          in.Address,
		Website:          in.Website,
		LogoURL:          in.LogoURL,
		RepresentativeID: in.RepresentativeID,
		Login:            in.Login,
		PasswordHash:     hash,
	}, now)
	if err != nil {
		return nil, err
	}
	return c, s.repo.Save(ctx, c)
}

func (s *companyService) Approve(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	c, err := s.repo.ByID(ctx, id)
	if err != nil {
		return nil, err
	}
	now := s.clock.Now()
	c2, err := c.Approve(now)
	if err != nil {
		return nil, err
	}
	return c2, s.repo.Save(ctx, c2)
}

func (s *companyService) Login(ctx context.Context, login, password string) (*domain.Company, error) {
	c, err := s.repo.ByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	if err := s.pass.Compare(c.Immutable().PasswordHash, password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	if !c.Immutable().Approved {
		return nil, fmt.Errorf("company not approved")
	}
	return c, nil
}

func (s *companyService) UpdateProfile(ctx context.Context, id uuid.UUID, in port.UpdateCompanyInput) (*domain.Company, error) {
	c, err := s.repo.ByID(ctx, id)
	if err != nil {
		return nil, err
	}
	now := s.clock.Now()
	c2, err := c.UpdateProfile(in.Title, in.Description, in.Contacts, in.Address, in.Website, in.LogoURL, now)
	if err != nil {
		return nil, err
	}
	return c2, s.repo.Save(ctx, c2)
}

func (s *companyService) ChangeCredentials(ctx context.Context, id uuid.UUID, login, password string) (*domain.Company, error) {
	c, err := s.repo.ByID(ctx, id)
	if err != nil {
		return nil, err
	}
	var hash string
	if password != "" {
		h, err := s.pass.Hash(password)
		if err != nil {
			return nil, err
		}
		hash = h
	}
	now := s.clock.Now()
	c2, err := c.ChangeCredentials(login, hash, now)
	if err != nil {
		return nil, err
	}
	return c2, s.repo.Save(ctx, c2)
}

func (s *companyService) Get(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	return s.repo.ByID(ctx, id)
}
