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

type companyRepo struct {
    q *pgqueries.Queries
}

func NewCompanyRepo(q *pgqueries.Queries) *companyRepo {
    return &companyRepo{q}
}

func (r *companyRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
    cdb, err := r.q.GetCompanyByID(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, domain.ErrNotFound
        }
        return nil, fmt.Errorf("error getting company by id: %w", err)
    }

    c, err := domain.ReconstructCompany(domain.CompanyImmutable{
        ID:              cdb.ID,
        Title:           cdb.Title,
        Description:     cdb.Description,
        Contacts:        cdb.Contacts,
        INN:             cdb.Inn,
        Address:         cdb.Address,
        Approved:        cdb.Approved,
        RepresentativeID: cdb.RepresentativeID,
        Login:           cdb.Login,
        PasswordHash:    cdb.PasswordHash,
        CreatedAt:       cdb.CreatedAt,
        UpdatedAt:       cdb.UpdatedAt,
    })
    if err != nil {
        return nil, fmt.Errorf("error reconstructing company: %w", err)
    }

    return c, nil
}

func (r *companyRepo) Save(ctx context.Context, c *domain.Company) error {
    _, err := r.GetByID(ctx, c.Immutable().ID)
    if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            return r.create(ctx, c)
        }
        return fmt.Errorf("error getting company by id: %w", err)
    }
    return r.update(ctx, c)
}

func (r *companyRepo) create(ctx context.Context, c *domain.Company) error {
    im := c.Immutable()
    err := r.q.CreateCompany(ctx, pgqueries.CreateCompanyParams{
        ID:              im.ID,
        Title:           im.Title,
        Description:     im.Description,
        Contacts:        im.Contacts,
        Inn:             im.INN,
        Address:         im.Address,
        Approved:        im.Approved,
        RepresentativeID: im.RepresentativeID,
        Login:           im.Login,
        PasswordHash:    im.PasswordHash,
        CreatedAt:       im.CreatedAt,
        UpdatedAt:       im.UpdatedAt,
    })
    if err != nil {
        if isUniqueViolationError(err) {
            return domain.ErrConflict
        }
        return fmt.Errorf("error creating company: %w", err)
    }
    return nil
}

func (r *companyRepo) update(ctx context.Context, c *domain.Company) error {
    im := c.Immutable()
    err := r.q.UpdateCompany(ctx, pgqueries.UpdateCompanyParams{
        ID:              im.ID,
        Title:           im.Title,
        Description:     im.Description,
        Contacts:        im.Contacts,
        Inn:             im.INN,
        Address:         im.Address,
        Approved:        im.Approved,
        RepresentativeID: im.RepresentativeID,
        Login:           im.Login,
        PasswordHash:    im.PasswordHash,
        CreatedAt:       im.CreatedAt,
        UpdatedAt:       im.UpdatedAt,
    })
    if err != nil {
        return fmt.Errorf("error updating company: %w", err)
    }
    return nil
}

// shared helper
