package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hr-platform-mosprom/internal/adapter/postgres/pgqueries"
	"github.com/hr-platform-mosprom/internal/core/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

type universityRepo struct {
	q *pgqueries.Queries
}

func NewUniversityRepo(q *pgqueries.Queries) *universityRepo {
	return &universityRepo{q}
}

func (r *universityRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.University, error) {
	universityFromDB, err := r.q.GetUniversityByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("error getting university by id: %w", err)
	}

	university, err := domain.ReconstructUniversity(domain.UniversityImmutable{
		ID:           universityFromDB.ID,
		Title:        universityFromDB.Title,
		Login:        universityFromDB.Login,
		PasswordHash: universityFromDB.PasswordHash,
		INN:          universityFromDB.Inn,
		Confirmed:    universityFromDB.Confirmed,
		Contacts:     universityFromDB.Contacts,
		CreatedAt:    universityFromDB.CreatedAt,
		UpdatedAt:    universityFromDB.UpdatedAt,
	})
	if err != nil {
		return nil, fmt.Errorf("error reconstructing university: %w", err)
	}

	return university, nil
}

func (r *universityRepo) GetByLogin(ctx context.Context, login string) (*domain.University, error) {
	universityFromDB, err := r.q.GetUniversityByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("error getting university by id: %w", err)
	}

	university, err := domain.ReconstructUniversity(domain.UniversityImmutable{
		ID:           universityFromDB.ID,
		Title:        universityFromDB.Title,
		Login:        universityFromDB.Login,
		PasswordHash: universityFromDB.PasswordHash,
		INN:          universityFromDB.Inn,
		Confirmed:    universityFromDB.Confirmed,
		Contacts:     universityFromDB.Contacts,
		CreatedAt:    universityFromDB.CreatedAt,
		UpdatedAt:    universityFromDB.UpdatedAt,
	})
	if err != nil {
		return nil, fmt.Errorf("error reconstructing university: %w", err)
	}

	return university, nil
}

func (r *universityRepo) Save(ctx context.Context, university *domain.University) error {
	_, err := r.GetByID(ctx, university.Immutable().ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			err = r.create(ctx, university)
			if err != nil {
				return fmt.Errorf("error creating university: %w", err)
			}
			return nil
		}
		return fmt.Errorf("error getting university by id: %w", err)
	}

	err = r.update(ctx, university)
	if err != nil {
		return fmt.Errorf("error updating university: %w", err)
	}

	return nil
}

func (r *universityRepo) create(ctx context.Context, university *domain.University) error {
	universityImmutable := university.Immutable()

	err := r.q.CreateUniversity(ctx, pgqueries.CreateUniversityParams{
		ID:           universityImmutable.ID,
		Title:        universityImmutable.Title,
		Login:        universityImmutable.Login,
		Inn:          universityImmutable.INN,
		Confirmed:    universityImmutable.Confirmed,
		PasswordHash: universityImmutable.PasswordHash,
		Contacts:     universityImmutable.Contacts,
		CreatedAt:    universityImmutable.CreatedAt,
		UpdatedAt:    universityImmutable.UpdatedAt,
	})
	if err != nil {
		if isUniqueViolationError(err) {
			return domain.ErrConflict
		}
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

func (r *universityRepo) update(ctx context.Context, university *domain.University) error {
	universityImmutable := university.Immutable()

	err := r.q.UpdateUniversity(ctx, pgqueries.UpdateUniversityParams{
		ID:           universityImmutable.ID,
		Title:        universityImmutable.Title,
		Login:        universityImmutable.Login,
		Inn:          universityImmutable.INN,
		Confirmed:    universityImmutable.Confirmed,
		PasswordHash: universityImmutable.PasswordHash,
		Contacts:     universityImmutable.Contacts,
		CreatedAt:    universityImmutable.CreatedAt,
		UpdatedAt:    universityImmutable.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}

func isUniqueViolationError(err error) bool {
	if err == nil {
		return false
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
