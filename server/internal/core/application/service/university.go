package service

import (
	"context"
	"fmt"

	"github.com/hr-platform-mosprom/internal/core/application/port"
	"github.com/hr-platform-mosprom/internal/core/domain"
)

type universityService struct {
	universityRepo  port.UniversityRepository
	passwordService port.PasswordService
	clock           port.Clock
}

func (s *universityService) SignUp(ctx context.Context, data port.SignUpUniversityData) error {
	err := s.validatePassword(data.Password)
	if err != nil {
		return fmt.Errorf("error validating password: %w", err)
	}

	passwordHash, err := s.passwordService.Hash(data.Password)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	university, err := domain.CreateUniversity(domain.CreateUniversityAttrs{
		Title:        data.Title,
		Login:        data.Login,
		PasswordHash: passwordHash,
		INN:          data.INN,
	}, s.clock.Now())
	if err != nil {
		return fmt.Errorf("error creating university: %w", err)
	}

	err = s.universityRepo.Save(ctx, university)
	if err != nil {
		return fmt.Errorf("error saving university: %w", err)
	}

	return nil
}

func (s *universityService) validatePassword(password string) error {
	if len(password) < 8 || len(password) > 64 {
		return fmt.Errorf("%w: invalid password length", domain.ErrInvariantViolated)
	}

	return nil
}
