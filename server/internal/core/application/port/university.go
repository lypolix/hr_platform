package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/hr-platform-mosprom/internal/core/domain"
)

type UniversityPreview struct {
	id    uuid.UUID
	title string
}

type UniversityRepository interface {
	Save(ctx context.Context, university *domain.University) error
	GetList(ctx context.Context, offset, limit int) []UniversityPreview
}

type SignUpUniversityData struct {
	Login    string
	Password string
	Title    string
	INN      string
}
