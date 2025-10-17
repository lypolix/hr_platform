package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hr-platform-mosprom/internal/core/application/port"
	"github.com/hr-platform-mosprom/internal/core/domain"
)

type vacancyService struct {
	vacRepo port.VacancyRepository
	coRepo  port.CompanyRepository  // в  дальнейшем для модерации вакансии
}

func NewVacancyService(vacRepo port.VacancyRepository, coRepo port.CompanyRepository) port.VacancyService {
	return &vacancyService{vacRepo: vacRepo, coRepo: coRepo}
}

func (s *vacancyService) Create(ctx context.Context, req domain.CreateVacancyRequest) (*domain.Vacancy, error) {
	co, err := s.coRepo.GetByID(ctx, req.CompanyID)
	if err != nil || co == nil {
		return nil, fmt.Errorf("company not found")
	}
	if !co.Approved {
		return nil, fmt.Errorf("company is not approved")
	}

	v := domain.NewVacancy(req)
	if err := s.vacRepo.Create(ctx, v); err != nil {
		return nil, fmt.Errorf("create vacancy: %w", err)
	}
	return v, nil
}

func (s *vacancyService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Vacancy, error) {
	return s.vacRepo.GetByID(ctx, id)
}

func (s *vacancyService) ListByCompany(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*domain.Vacancy, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.vacRepo.GetByCompanyID(ctx, companyID, limit, offset)
}

func (s *vacancyService) ListAll(ctx context.Context, limit, offset int) ([]*domain.Vacancy, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.vacRepo.GetAll(ctx, limit, offset)
}

func (s *vacancyService) Search(ctx context.Context, filter domain.VacancyFilter) ([]*domain.Vacancy, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	return s.vacRepo.GetWithFilters(ctx, filter)
}

func (s *vacancyService) Update(ctx context.Context, id uuid.UUID, req domain.UpdateVacancyRequest) (*domain.Vacancy, error) {
	v, err := s.vacRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get vacancy: %w", err)
	}
	v.Update(req)
	if err := s.vacRepo.Update(ctx, v); err != nil {
		return nil, fmt.Errorf("update vacancy: %w", err)
	}
	return v, nil
}

func (s *vacancyService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.vacRepo.Delete(ctx, id)
}

func (s *vacancyService) Activate(ctx context.Context, id uuid.UUID) error {
	v, err := s.vacRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	v.Activate()
	return s.vacRepo.Update(ctx, v)
}

func (s *vacancyService) Deactivate(ctx context.Context, id uuid.UUID) error {
	v, err := s.vacRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	v.Deactivate()
	return s.vacRepo.Update(ctx, v)
}
