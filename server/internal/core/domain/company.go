package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type (
	Company struct {
		id               uuid.UUID
		title            string
		description      string
		contacts         string
		inn              string
		address          string
		approved         bool
		representativeID uuid.UUID
		login            string
		passwordHash     string
		createdAt        time.Time
		updatedAt        time.Time
	}

	CompanyImmutable struct {
		ID               uuid.UUID
		Title            string
		Description      string
		Contacts         string
		INN              string
		Address          string
		LogoURL          string
		Approved         bool
		RepresentativeID uuid.UUID
		Login            string
		PasswordHash     string
		CreatedAt        time.Time
		UpdatedAt        time.Time
	}

	CreateCompanyAttrs struct {
		Title            string
		Description      string
		Contacts         string
		INN              string
		Address          string
		LogoURL          string
		RepresentativeID uuid.UUID
		Login            string
		PasswordHash     string
	}
)

func (c *Company) Immutable() CompanyImmutable {
	return CompanyImmutable{
		ID:               c.id,
		Title:            c.title,
		Description:      c.description,
		Contacts:         c.contacts,
		INN:              c.inn,
		Address:          c.address,
		Approved:         c.approved,
		RepresentativeID: c.representativeID,
		Login:            c.login,
		PasswordHash:     c.passwordHash,
		CreatedAt:        c.createdAt,
		UpdatedAt:        c.updatedAt,
	}
}

func (c *Company) checkInvariants() error {
	if c.id == uuid.Nil {
		return fmt.Errorf("%w: nil id", ErrInvariantViolated)
	}
	if l := len(c.title); l < 1 || l > 512 {
		return fmt.Errorf("%w: invalid title length", ErrInvariantViolated)
	}
	if len(c.inn) != 10 {
		return fmt.Errorf("%w: invalid INN length", ErrInvariantViolated)
	}
	if l := len(c.login); l < 4 || l > 128 {
		return fmt.Errorf("%w: invalid login length", ErrInvariantViolated)
	}
	if len(c.passwordHash) < 10 {
		return fmt.Errorf("%w: weak password hash", ErrInvariantViolated)
	}
	if c.createdAt.IsZero() {
		return fmt.Errorf("%w: zero creation time", ErrInvariantViolated)
	}
	if c.updatedAt.IsZero() {
		return fmt.Errorf("%w: zero updation time", ErrInvariantViolated)
	}
	if c.createdAt.After(c.updatedAt) {
		return fmt.Errorf("%w: creation time is after updation time", ErrInvariantViolated)
	}
	return nil
}

func CreateCompany(attrs CreateCompanyAttrs, at time.Time) (*Company, error) {
	imm := CompanyImmutable{
		ID:               uuid.New(),
		Title:            attrs.Title,
		Description:      attrs.Description,
		Contacts:         attrs.Contacts,
		INN:              attrs.INN,
		Address:          attrs.Address,
		LogoURL:          attrs.LogoURL,
		Approved:         false,
		RepresentativeID: attrs.RepresentativeID,
		Login:            attrs.Login,
		PasswordHash:     attrs.PasswordHash,
		CreatedAt:        at,
		UpdatedAt:        at,
	}
	return ReconstructCompany(imm)
}

func ReconstructCompany(immutable CompanyImmutable) (*Company, error) {
	c := &Company{
		id:               immutable.ID,
		title:            immutable.Title,
		description:      immutable.Description,
		contacts:         immutable.Contacts,
		inn:              immutable.INN,
		address:          immutable.Address,
		
		approved:         immutable.Approved,
		representativeID: immutable.RepresentativeID,
		login:            immutable.Login,
		passwordHash:     immutable.PasswordHash,
		createdAt:        immutable.CreatedAt,
		updatedAt:        immutable.UpdatedAt,
	}
	return c, c.checkInvariants()
}

// Мутации через Immutable + Reconstruct
func (c *Company) Approve(at time.Time) (*Company, error) {
	imm := c.Immutable()
	imm.Approved = true
	imm.UpdatedAt = at
	return ReconstructCompany(imm)
}

func (c *Company) UpdateProfile(title, description, contacts, address string, at time.Time) (*Company, error) {
	imm := c.Immutable()
	if title != "" {
		imm.Title = title
	}
	if description != "" {
		imm.Description = description
	}
	if contacts != "" {
		imm.Contacts = contacts
	}
	if address != "" {
		imm.Address = address
	}
	
	imm.UpdatedAt = at
	return ReconstructCompany(imm)
}

func (c *Company) ChangeCredentials(login, passwordHash string, at time.Time) (*Company, error) {
	imm := c.Immutable()
	if login != "" {
		imm.Login = login
	}
	if passwordHash != "" {
		imm.PasswordHash = passwordHash
	}
	imm.UpdatedAt = at
	return ReconstructCompany(imm)
}
