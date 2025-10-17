package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hr-platform-mosprom/internal/core/application/port"
	"github.com/hr-platform-mosprom/internal/core/domain"
)

type companyService struct {
	companyRepo     port.CompanyRepository
	passwordService port.PasswordService
	tokenService    port.TokenService
	clock           port.Clock
}

type CompanyServiceDeps struct {
	CompanyRepo     port.CompanyRepository
	PasswordService port.PasswordService
	TokenService    port.TokenService
	Clock           port.Clock
}

func NewCompanyService(d CompanyServiceDeps) *companyService {
	return &companyService{
		companyRepo:     d.CompanyRepo,
		passwordService: d.PasswordService,
		tokenService:    d.TokenService,
		clock:           d.Clock,
	}
}

func (s *companyService) SignUp(ctx context.Context, data port.SignUpCompanyData) (*port.CompanyWithTokenResult, error) {
	if err := s.validatePassword(data.Password); err != nil {
		return nil, fmt.Errorf("error validating password: %w", err)
	}

	passwordHash, err := s.passwordService.Hash(data.Password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	company, err := domain.CreateCompany(domain.CreateCompanyAttrs{
		Title:            data.Title,
		Description:      data.Description,
		Contacts:         data.Contacts,
		INN:              data.INN,
		Address:          data.Address,
		Website:          data.Website,
		LogoURL:          data.LogoURL,
		RepresentativeID: data.RepresentativeID,
		Login:            data.Login,
		PasswordHash:     passwordHash,
	}, s.clock.Now())
	if err != nil {
		return nil, fmt.Errorf("error creating company: %w", err)
	}

	if err := s.companyRepo.Save(ctx, company); err != nil {
		return nil, fmt.Errorf("error saving company: %w", err)
	}

	ci := company.Immutable()
	token, err := s.tokenService.Generate(port.TokenPayload{
		Sub:  ci.ID,
		Role: port.RoleCompany,
	})
	if err != nil {
		return nil, fmt.Errorf("error generating token: %w", err)
	}

	return &port.CompanyWithTokenResult{
		CompanyResult: port.CompanyResult{
			Login: ci.Login,
			Title: ci.Title,
			INN:   ci.INN,
		},
		Token: token,
	}, nil
}

func (s *companyService) SignIn(ctx context.Context, data port.SignInCompanyData) (*port.CompanyWithTokenResult, error) {
	company, err := s.companyRepo.GetByLogin(ctx, data.Login)
	if err != nil {
		return nil, fmt.Errorf("error getting company by login: %w", err)
	}

	ci := company.Immutable()
	if !s.passwordService.Check(data.Password, ci.PasswordHash) {
		return nil, domain.ErrUnauthorized
	}
	if !ci.Approved {
		return nil, fmt.Errorf("company is not approved")
	}

	token, err := s.tokenService.Generate(port.TokenPayload{
		Sub:  ci.ID,
		Role: port.RoleCompany,
	})
	if err != nil {
		return nil, fmt.Errorf("error generating token: %w", err)
	}

	return &port.CompanyWithTokenResult{
		CompanyResult: port.CompanyResult{
			Login: ci.Login,
			Title: ci.Title,
			INN:   ci.INN,
		},
		Token: token,
	}, nil
}

func (s *companyService) Approve(ctx context.Context, actor domain.Actor, companyID uuid.UUID) error {
	if actor.Role != domain.RoleAdmin {
		return fmt.Errorf("admin role required: %w", domain.ErrForbidden)
	}

	company, err := s.companyRepo.GetByID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("error getting company by id: %w", err)
	}

	company2, err := company.Approve(s.clock.Now())
	if err != nil {
		return fmt.Errorf("error approving company: %w", err)
	}

	return s.companyRepo.Save(ctx, company2)
}


func (s *companyService) UpdateProfile(ctx context.Context, actor domain.Actor, data port.UpdateCompanyProfileData) (*port.CompanyResult, error) {
	if actor.Role != domain.RoleCompany {
		return nil, fmt.Errorf("company role required: %w", domain.ErrForbidden)
	}

	company, err := s.companyRepo.GetByID(ctx, actor.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting company by id: %w", err)
	}

	company2, err := company.UpdateProfile(
		data.Title, data.Description, data.Contacts, data.Address, data.Website, data.LogoURL,
		s.clock.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("error updating profile: %w", err)
	}

	if err := s.companyRepo.Save(ctx, company2); err != nil {
		return nil, fmt.Errorf("error saving company: %w", err)
	}

	ci := company2.Immutable()
	return &port.CompanyResult{
		Login: ci.Login,
		Title: ci.Title,
		INN:   ci.INN,
	}, nil
}

func (s *companyService) ChangeCredentials(ctx context.Context, actor domain.Actor, data port.ChangeCompanyCredentialsData) error {
	if actor.Role != domain.RoleCompany {
		return fmt.Errorf("company role required: %w", domain.ErrForbidden)
	}

	company, err := s.companyRepo.GetByID(ctx, actor.ID)
	if err != nil {
		return fmt.Errorf("error getting company by id: %w", err)
	}

	var newHash string
	if data.Password != "" {
		if err := s.validatePassword(data.Password); err != nil {
			return fmt.Errorf("error validating password: %w", err)
		}
		h, err := s.passwordService.Hash(data.Password)
		if err != nil {
			return fmt.Errorf("error hashing password: %w", err)
		}
		newHash = h
	}

	company2, err := company.ChangeCredentials(data.Login, newHash, s.clock.Now())
	if err != nil {
		return fmt.Errorf("error changing credentials: %w", err)
	}

	return s.companyRepo.Save(ctx, company2)
}

func (s *companyService) validatePassword(password string) error {
	if len(password) < 8 || len(password) > 64 {
		return fmt.Errorf("%w: invalid password length", domain.ErrInvariantViolated)
	}
	return nil
}
