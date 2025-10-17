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
	tokenService    port.TokenService
	clock           port.Clock
}

func (s *universityService) SignUp(ctx context.Context, data port.SignUpUniversityData) (*port.UniversityWithTokenResult, error) {
	err := s.validatePassword(data.Password)
	if err != nil {
		return nil, fmt.Errorf("error validating password: %w", err)
	}

	passwordHash, err := s.passwordService.Hash(data.Password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	university, err := domain.CreateUniversity(domain.CreateUniversityAttrs{
		Title:        data.Title,
		Login:        data.Login,
		PasswordHash: passwordHash,
		INN:          data.INN,
	}, s.clock.Now())
	if err != nil {
		return nil, fmt.Errorf("error creating university: %w", err)
	}

	err = s.universityRepo.Save(ctx, university)
	if err != nil {
		return nil, fmt.Errorf("error saving university: %w", err)
	}

	universityImmutable := university.Immutable()

	token, err := s.tokenService.Generate(port.TokenPayload{
		Sub:  universityImmutable.ID,
		Role: port.RoleUniversity,
	})
	if err != nil {
		return nil, fmt.Errorf("error generating token: %w", err)
	}

	universityWithTokenResult := &port.UniversityWithTokenResult{
		UniversityResult: port.UniversityResult{
			Login: universityImmutable.Login,
			Title: universityImmutable.Title,
			INN:   universityImmutable.INN,
		},
		Token: token,
	}

	return universityWithTokenResult, nil
}

func (s *universityService) SignIn(ctx context.Context, data port.SignInUniversityData) (*port.UniversityWithTokenResult, error) {
	university, err := s.universityRepo.GetByLogin(ctx, data.Login)
	if err != nil {
		return nil, fmt.Errorf("error getting university by login: %w", err)
	}

	universityImmutable := university.Immutable()

	if !s.passwordService.Check(data.Password, universityImmutable.PasswordHash) {
		return nil, domain.ErrUnauthorized
	}

	token, err := s.tokenService.Generate(port.TokenPayload{
		Sub:  universityImmutable.ID,
		Role: port.RoleUniversity,
	})
	if err != nil {
		return nil, fmt.Errorf("error generating token: %w", err)
	}

	universityWithTokenResult := &port.UniversityWithTokenResult{
		UniversityResult: port.UniversityResult{
			Login: universityImmutable.Login,
			Title: universityImmutable.Title,
			INN:   universityImmutable.INN,
		},
		Token: token,
	}

	return universityWithTokenResult, nil
}

func (s *universityService) ChangePassword(ctx context.Context, actor domain.Actor, data port.ChangeUniversityPasswordData) error {
	if actor.Role != domain.RoleUniversity {
		return fmt.Errorf("university role required: %w", domain.ErrForbidden)
	}

	university, err := s.universityRepo.GetByID(ctx, actor.ID)
	if err != nil {
		return fmt.Errorf("error getting university by id: %w", err)
	}

	if !s.passwordService.Check(data.CurrentPassword, university.Immutable().PasswordHash) {
		return domain.ErrUnauthorized
	}

	err = s.validatePassword(data.NewPassword)
	if err != nil {
		return fmt.Errorf("error validating password: %w", err)
	}

	newPasswordHash, err := s.passwordService.Hash(data.NewPassword)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	err = university.SetPasswordHash(newPasswordHash, s.clock.Now())
	if err != nil {
		return fmt.Errorf("error setting password hash: %w", err)
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
