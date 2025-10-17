package domain

import (
	"time"

	"github.com/google/uuid"
)

type Company struct {
	ID               uuid.UUID `json:"id" db:"id"`
	Title            string    `json:"title" db:"title"`                  // название компании
	INN              string    `json:"inn" db:"inn"`                      // ИНН
	Description      string    `json:"description" db:"description"`      // описание
	Contacts         string    `json:"contacts" db:"contacts"`            // контакты в свободной форме
	Address          string    `json:"address" db:"address"`              // адрес (опц.)
	Website          string    `json:"website" db:"website"`              // сайт (опц.)
	LogoURL          string    `json:"logo_url" db:"logo_url"`            // логотип (опц.)
	Approved         bool      `json:"approved" db:"approved"`            // прошла модерацию?
	RepresentativeID uuid.UUID `json:"representative_id" db:"representative_id"`

	// Данные для авторизации представителя компании
	Login        string `json:"login" db:"login"`
	PasswordHash string `json:"-" db:"password_hash"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateCompanyRequest struct {
	Title            string    `json:"title" binding:"required"`
	INN              string    `json:"inn" binding:"required,len=10"`
	Description      string    `json:"description" binding:"required"`
	Contacts         string    `json:"contacts" binding:"required"`
	Address          string    `json:"address"`
	Website          string    `json:"website"`
	LogoURL          string    `json:"logo_url"`
	RepresentativeID uuid.UUID `json:"representative_id" binding:"required"`

	Login    string `json:"login" binding:"required,min=4"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateCompanyRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Contacts    string `json:"contacts"`
	Address     string `json:"address"`
	Website     string `json:"website"`
	LogoURL     string `json:"logo_url"`
}

func NewCompany(req CreateCompanyRequest, passwordHash string) *Company {
	now := time.Now()
	return &Company{
		ID:               uuid.New(),
		Title:            req.Title,
		INN:              req.INN,
		Description:      req.Description,
		Contacts:         req.Contacts,
		Address:          req.Address,
		Website:          req.Website,
		LogoURL:          req.LogoURL,
		Approved:         false,
		RepresentativeID: req.RepresentativeID,
		Login:            req.Login,
		PasswordHash:     passwordHash,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

func (c *Company) Update(req UpdateCompanyRequest) {
	if req.Title != "" {
		c.Title = req.Title
	}
	if req.Description != "" {
		c.Description = req.Description
	}
	if req.Contacts != "" {
		c.Contacts = req.Contacts
	}
	if req.Address != "" {
		c.Address = req.Address
	}
	if req.Website != "" {
		c.Website = req.Website
	}
	if req.LogoURL != "" {
		c.LogoURL = req.LogoURL
	}
	c.UpdatedAt = time.Now()
}

func (c *Company) Approve() {
	c.Approved = true
	c.UpdatedAt = time.Now()
}
