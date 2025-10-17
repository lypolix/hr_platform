package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type (
	University struct {
		id           uuid.UUID
		title        string
		login        string
		passwordHash string
		inn          string
		confirmed    bool
		createdAt    time.Time
		updatedAt    time.Time
	}

	UniversityImmutable struct {
		ID           uuid.UUID
		Title        string
		Login        string
		PasswordHash string
		INN          string
		Confirmed    bool
		CreatedAt    time.Time
		UpdatedAt    time.Time
	}

	CreateUniversityAttrs struct {
		Title        string
		Login        string
		PasswordHash string
		INN          string
	}
)

func (u *University) Immutable() UniversityImmutable {
	return UniversityImmutable{
		ID:           u.id,
		Title:        u.title,
		Login:        u.login,
		PasswordHash: u.passwordHash,
		INN:          u.inn,
		Confirmed:    u.confirmed,
	}
}

func (u *University) SetPasswordHash(newPasswordHash string, at time.Time) error {
	u.passwordHash = newPasswordHash
	u.updatedAt = at

	return u.checkInvariants()
}

func (u *University) checkInvariants() error {
	if u.id == uuid.Nil {
		return fmt.Errorf("%w: nil id", ErrInvariantViolated)
	}

	if len(u.title) < 1 || len(u.title) > 512 {
		return fmt.Errorf("%w: invalid title length", ErrInvariantViolated)
	}

	if u.createdAt.IsZero() {
		return fmt.Errorf("%w: zero creation time", ErrInvariantViolated)
	}

	if u.updatedAt.IsZero() {
		return fmt.Errorf("%w: zero updation time", ErrInvariantViolated)
	}

	if u.createdAt.After(u.updatedAt) {
		return fmt.Errorf("%w: creation time is after updation time", ErrInvariantViolated)
	}

	return nil
}

func CreateUniversity(attrs CreateUniversityAttrs, at time.Time) (*University, error) {
	u := &University{
		id:           uuid.New(),
		title:        attrs.Title,
		login:        attrs.Login,
		passwordHash: attrs.PasswordHash,
		inn:          attrs.INN,
		confirmed:    false,
		createdAt:    at,
		updatedAt:    at,
	}

	return u, u.checkInvariants()
}

func ReconstructUniversity(immutable UniversityImmutable) (*University, error) {
	u := &University{
		id:           immutable.ID,
		title:        immutable.Title,
		login:        immutable.Login,
		passwordHash: immutable.PasswordHash,
		inn:          immutable.INN,
		confirmed:    immutable.Confirmed,
		createdAt:    immutable.CreatedAt,
		updatedAt:    immutable.UpdatedAt,
	}

	return u, u.checkInvariants()
}
