package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hr-platform-mosprom/internal/core/application/port"
	"github.com/hr-platform-mosprom/internal/core/domain"
)

type vacancyService struct {
	repo    port.VacancyRepository
	company port.CompanyRepository
	clock   port.Clock
}

// БЫЛО: func NewVacancyService(r port.VacancyRepository, cr port.CompanyRepository, c port.Clock) *v vacancyService {
// ДОЛЖНО БЫТЬ:
func NewVacancyService(r port.VacancyRepository, cr port.CompanyRepository, c port.Clock) *vacancyService {
	return &vacancyService{repo: r, company: cr, clock: c}
}

func (s *vacancyService) Create(ctx context.Context, in port.CreateVacancyInput) (*domain.Vacancy, error) {
	co, err := s.company.ByID(ctx, in.CompanyID)
	if err != nil || co == nil {
		return nil, fmt.Errorf("company not found")
	}
	if !co.Immutable().Approved {
		return nil, fmt.Errorf("company not approved")
	}

	now := s.clock.Now()

	v, err := domain.CreateVacancy(domain.CreateVacancyAttrs{
		CompanyID:        in.CompanyID,
		Title:            in.Title,
		Description:      in.Description,
		Contacts:         in.Contacts,
		Requirements:     in.Requirements,
		Responsibilities: in.Responsibilities,
		Conditions:       in.Conditions,
		SalaryFrom:       in.SalaryFrom,
		SalaryTo:         in.SalaryTo,
		Employment:       in.Employment,
		Schedule:         in.Schedule,
		Experience:       in.Experience,
		Education:        in.Education,
		Location:         in.Location,
	}, now)
	if err != nil {
		return nil, err
	}
	return v, s.repo.Save(ctx, v)
}

func (s *vacancyService) Update(ctx context.Context, id uuid.UUID, in port.UpdateVacancyInput) (*domain.Vacancy, error) {
	v, err := s.repo.ByID(ctx, id)
	if err != nil {
		return nil, err
	}

	now := s.clock.Now()

	v2, err := v.Update(domain.CreateVacancyAttrs{
		Title:            in.Title,
		Description:      in.Description,
		Contacts:         in.Contacts,
		Requirements:     in.Requirements,
		Responsibilities: in.Responsibilities,
		Conditions:       in.Conditions,
		SalaryFrom:       in.SalaryFrom,
		SalaryTo:         in.SalaryTo,
		Employment:       in.Employment,
		Schedule:         in.Schedule,
		Experience:       in.Experience,
		Education:        in.Education,
		Location:         in.Location,
	}, now)
	if err != nil {
		return nil, err
	}

	if in.IsActive != nil {
		if *in.IsActive {
			v2, err = v2.Activate(now)
		} else {
			v2, err = v2.Deactivate(now)
		}
		if err != nil {
			return nil, err
		}
	}

	return v2, s.repo.Save(ctx, v2)
}

func (s *vacancyService) Get(ctx context.Context, id uuid.UUID) (*domain.Vacancy, error) {
	return s.repo.ByID(ctx, id)
}

func (s *vacancyService) ListByCompany(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*domain.Vacancy, error) {
	return s.repo.ByCompany(ctx, companyID, limit, offset)
}

func (s *vacancyService) ListAll(ctx context.Context, limit, offset int) ([]*domain.Vacancy, error) {
	return s.repo.All(ctx, limit, offset)
}

func (s *vacancyService) Search(ctx context.Context, f domain.VacancyFilter) ([]*domain.Vacancy, error) {
	return s.repo.Search(ctx, f)
}

func (s *vacancyService) Activate(ctx context.Context, id uuid.UUID) (*domain.Vacancy, error) {
	v, err := s.repo.ByID(ctx, id)
	if err != nil {
		return nil, err
	}
	now := s.clock.Now()
	v2, err := v.Activate(now)
	if err != nil {
		return nil, err
	}
	return v2, s.repo.Save(ctx, v2)
}

func (s *vacancyService) Deactivate(ctx context.Context, id uuid.UUID) (*domain.Vacancy, error) {
	v, err := s.repo.ByID(ctx, id)
	if err != nil {
		return nil, err
	}
	now := s.clock.Now()
	v2, err := v.Deactivate(now)
	if err != nil {
		return nil, err
	}
	return v2, s.repo.Save(ctx, v2)
}

func (s *vacancyService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
