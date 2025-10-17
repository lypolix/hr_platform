package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/hr-platform-mosprom/internal/core/domain"
)

type UniversityRepository interface {
	Save(ctx context.Context, university *domain.University) error
	GetByLogin(ctx context.Context, login string) (*domain.University, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.University, error)
}

type SignUpUniversityData struct {
	Login    string
	Password string
	Title    string
	INN      string
}

type SignInUniversityData struct {
	Login    string
	Password string
}

type UniversityResult struct {
	Login string
	Title string
	INN   string
}

type UniversityWithTokenResult struct {
	UniversityResult
	Token string
}

type ChangeUniversityPasswordData struct {
	NewPassword     string
	CurrentPassword string
}
