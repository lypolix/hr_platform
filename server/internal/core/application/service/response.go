package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/hr-platform-mosprom/internal/core/application/port"
	"github.com/hr-platform-mosprom/internal/core/domain"
)

type responseService struct {
	repo  port.ResponseRepository
	vRepo port.VacancyRepository
	clock port.Clock
}

func NewResponseService(r port.ResponseRepository, v port.VacancyRepository, c port.Clock) *responseService {
	return &responseService{repo: r, vRepo: v, clock: c}
}

func (s *responseService) Create(ctx context.Context, in port.CreateResponseInput) (*domain.Response, error) {
	if _, err := s.vRepo.ByID(ctx, in.VacancyID); err != nil {
		return nil, err
	}
	now := s.clock.Now()
	r, err := domain.CreateResponse(domain.CreateResponseAttrs{
		VacancyID:   in.VacancyID,
		FullName:    in.FullName,
		Email:       in.Email,
		Phone:       in.Phone,
		CoverLetter: in.CoverLetter,
		ResumeURL:   in.ResumeURL,
	}, now)
	if err != nil {
		return nil, err
	}
	return r, s.repo.Save(ctx, r)
}

func (s *responseService) SetStatus(ctx context.Context, id uuid.UUID, status string) (*domain.Response, error) {
	r, err := s.repo.ByID(ctx, id)
	if err != nil {
		return nil, err
	}
	now := s.clock.Now()
	r2, err := r.SetStatus(status, now)
	if err != nil {
		return nil, err
	}
	return r2, s.repo.Save(ctx, r2)
}
