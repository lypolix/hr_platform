package postgres

import (
    "context"
    "database/sql"
    "errors"
    "fmt"

    "github.com/google/uuid"
    "github.com/hr-platform-mosprom/internal/adapter/postgres/pgqueries"
    "github.com/hr-platform-mosprom/internal/core/domain"
)

type vacancyRepo struct {
    q *pgqueries.Queries
}

func NewVacancyRepo(q *pgqueries.Queries) *vacancyRepo {
    return &vacancyRepo{q}
}

func (r *vacancyRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Vacancy, error) {
    vdb, err := r.q.GetVacancyByID(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, domain.ErrNotFound
        }
        return nil, fmt.Errorf("error getting vacancy by id: %w", err)
    }

    v, err := domain.ReconstructVacancy(domain.VacancyImmutable{
        ID:              vdb.ID,
        CompanyID:       vdb.CompanyID,
        Title:           vdb.Title,
        Description:     vdb.Description,
        Contacts:        vdb.Contacts,
        Requirements:    vdb.Requirements,
        Responsibilities: vdb.Responsibilities,
        Conditions:      vdb.Conditions,
        Employment:      vdb.Employment,
        Schedule:        vdb.Schedule,
        Experience:      vdb.Experience,
        Education:       vdb.Education,
        Location:        vdb.Location,
        IsActive:        vdb.IsActive,
        CreatedAt:       vdb.CreatedAt,
        UpdatedAt:       vdb.UpdatedAt,
    })
    if err != nil {
        return nil, fmt.Errorf("error reconstructing vacancy: %w", err)
    }

    return v, nil
}

func (r *vacancyRepo) Save(ctx context.Context, v *domain.Vacancy) error {
    _, err := r.GetByID(ctx, v.Immutable().ID)
    if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            return r.create(ctx, v)
        }
        return fmt.Errorf("error getting vacancy by id: %w", err)
    }
    return r.update(ctx, v)
}

func (r *vacancyRepo) create(ctx context.Context, v *domain.Vacancy) error {
    im := v.Immutable()
    err := r.q.CreateVacancy(ctx, pgqueries.CreateVacancyParams{
        ID:              im.ID,
        CompanyID:       im.CompanyID,
        Title:           im.Title,
        Description:     im.Description,
        Contacts:        im.Contacts,
        Requirements:    im.Requirements,
        Responsibilities: im.Responsibilities,
        Conditions:      im.Conditions,
        Employment:      im.Employment,
        Schedule:        im.Schedule,
        Experience:      im.Experience,
        Education:       im.Education,
        Location:        im.Location,
        IsActive:        im.IsActive,
        CreatedAt:       im.CreatedAt,
        UpdatedAt:       im.UpdatedAt,
    })
    if err != nil {
        return fmt.Errorf("error creating vacancy: %w", err)
    }
    return nil
}

func (r *vacancyRepo) update(ctx context.Context, v *domain.Vacancy) error {
    im := v.Immutable()
    err := r.q.UpdateVacancy(ctx, pgqueries.UpdateVacancyParams{
        ID:              im.ID,
        CompanyID:       im.CompanyID,
        Title:           im.Title,
        Description:     im.Description,
        Contacts:        im.Contacts,
        Requirements:    im.Requirements,
        Responsibilities: im.Responsibilities,
        Conditions:      im.Conditions,
        Employment:      im.Employment,
        Schedule:        im.Schedule,
        Experience:      im.Experience,
        Education:       im.Education,
        Location:        im.Location,
        IsActive:        im.IsActive,
        CreatedAt:       im.CreatedAt,
        UpdatedAt:       im.UpdatedAt,
    })
    if err != nil {
        return fmt.Errorf("error updating vacancy: %w", err)
    }
    return nil
}
