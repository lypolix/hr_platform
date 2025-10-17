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

type responseRepo struct {
    q *pgqueries.Queries
}

func NewResponseRepo(q *pgqueries.Queries) *responseRepo {
    return &responseRepo{q}
}

func (r *responseRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Response, error) {
    rdb, err := r.q.GetResponseByID(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, domain.ErrNotFound
        }
        return nil, fmt.Errorf("error getting response by id: %w", err)
    }

    resp, err := domain.ReconstructResponse(domain.ResponseImmutable{
        ID:         rdb.ID,
        VacancyID:  rdb.VacancyID,
        FullName:   rdb.FullName,
        Email:      rdb.Email,
        Phone:      rdb.Phone,
        CoverLetter: rdb.CoverLetter,
        ResumeURL:  rdb.ResumeUrl,
        Status:     rdb.Status,
        CreatedAt:  rdb.CreatedAt,
        UpdatedAt:  rdb.UpdatedAt,
    })
    if err != nil {
        return nil, fmt.Errorf("error reconstructing response: %w", err)
    }

    return resp, nil
}

func (r *responseRepo) Save(ctx context.Context, resp *domain.Response) error {
    _, err := r.GetByID(ctx, resp.Immutable().ID)
    if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            return r.create(ctx, resp)
        }
        return fmt.Errorf("error getting response by id: %w", err)
    }
    return r.update(ctx, resp)
}

func (r *responseRepo) create(ctx context.Context, resp *domain.Response) error {
    im := resp.Immutable()
    err := r.q.CreateResponse(ctx, pgqueries.CreateResponseParams{
        ID:          im.ID,
        VacancyID:   im.VacancyID,
        FullName:    im.FullName,
        Email:       im.Email,
        Phone:       im.Phone,
        CoverLetter: im.CoverLetter,
        ResumeUrl:   im.ResumeURL,
        Status:      im.Status,
        CreatedAt:   im.CreatedAt,
        UpdatedAt:   im.UpdatedAt,
    })
    if err != nil {
        return fmt.Errorf("error creating response: %w", err)
    }
    return nil
}

func (r *responseRepo) update(ctx context.Context, resp *domain.Response) error {
    im := resp.Immutable()
    // Часто обновляется только статус; если у тебя общий Update — оставим полный апдейт
    err := r.q.UpdateResponseStatus(ctx, pgqueries.UpdateResponseStatusParams{
        ID:        im.ID,
        Status:    im.Status,
        UpdatedAt: im.UpdatedAt,
    })
    if err != nil {
        return fmt.Errorf("error updating response: %w", err)
    }
    return nil
}
