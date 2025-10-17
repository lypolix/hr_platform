package domain

import (
	"time"

	"github.com/google/uuid"
)

type Vacancy struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CompanyID   uuid.UUID `json:"company_id" db:"company_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Contacts    string    `json:"contacts" db:"contacts"`

	Requirements     string `json:"requirements" db:"requirements"`
	Responsibilities string `json:"responsibilities" db:"responsibilities"`
	Conditions       string `json:"conditions" db:"conditions"`

	SalaryFrom *int `json:"salary_from" db:"salary_from"`
	SalaryTo   *int `json:"salary_to" db:"salary_to"`

	Employment string `json:"employment" db:"employment"` // занятость 
	Schedule   string `json:"schedule" db:"schedule"`     //  онлайн, гибрид и тд
	Experience string `json:"experience" db:"experience"` //0/1-3/3-6/6+
	Education  string `json:"education" db:"education"`   // -/secondary/higher
	Location   string `json:"location" db:"location"`

	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateVacancyRequest struct {
	CompanyID       uuid.UUID `json:"company_id" binding:"required"`
	Title           string    `json:"title" binding:"required"`
	Description     string    `json:"description" binding:"required"`
	Contacts        string    `json:"contacts" binding:"required"`
	Requirements     string `json:"requirements" binding:"required"`
	Responsibilities string `json:"responsibilities" binding:"required"`
	Conditions       string `json:"conditions"`
	SalaryFrom       *int   `json:"salary_from"`
	SalaryTo         *int   `json:"salary_to"`
	Employment       string `json:"employment" binding:"required"`
	Schedule         string `json:"schedule" binding:"required"`
	Experience       string `json:"experience" binding:"required"`
	Education        string `json:"education" binding:"required"`
	Location         string `json:"location" binding:"required"`
}

type UpdateVacancyRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Contacts        string `json:"contacts"`
	Requirements     string `json:"requirements"`
	Responsibilities string `json:"responsibilities"`
	Conditions       string `json:"conditions"`
	SalaryFrom       *int   `json:"salary_from"`
	SalaryTo         *int   `json:"salary_to"`
	Employment       string `json:"employment"`
	Schedule         string `json:"schedule"`
	Experience       string `json:"experience"`
	Education        string `json:"education"`
	Location         string `json:"location"`
	IsActive         *bool  `json:"is_active"`
}

type VacancyFilter struct {
	CompanyID  *uuid.UUID
	Location   *string
	Employment *string
	Schedule   *string
	Experience *string
	Education  *string
	IsActive   *bool
	Limit      int
	Offset     int
}

func NewVacancy(req CreateVacancyRequest) *Vacancy {
	now := time.Now()
	return &Vacancy{
		ID:              uuid.New(),
		CompanyID:       req.CompanyID,
		Title:           req.Title,
		Description:     req.Description,
		Contacts:        req.Contacts,
		Requirements:    req.Requirements,
		Responsibilities: req.Responsibilities,
		Conditions:      req.Conditions,
		SalaryFrom:      req.SalaryFrom,
		SalaryTo:        req.SalaryTo,
		Employment:      req.Employment,
		Schedule:        req.Schedule,
		Experience:      req.Experience,
		Education:       req.Education,
		Location:        req.Location,
		IsActive:        true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func (v *Vacancy) Update(req UpdateVacancyRequest) {
	if req.Title != "" {
		v.Title = req.Title
	}
	if req.Description != "" {
		v.Description = req.Description
	}
	if req.Contacts != "" {
		v.Contacts = req.Contacts
	}
	if req.Requirements != "" {
		v.Requirements = req.Requirements
	}
	if req.Responsibilities != "" {
		v.Responsibilities = req.Responsibilities
	}
	if req.Conditions != "" {
		v.Conditions = req.Conditions
	}
	if req.SalaryFrom != nil {
		v.SalaryFrom = req.SalaryFrom
	}
	if req.SalaryTo != nil {
		v.SalaryTo = req.SalaryTo
	}
	if req.Employment != "" {
		v.Employment = req.Employment
	}
	if req.Schedule != "" {
		v.Schedule = req.Schedule
	}
	if req.Experience != "" {
		v.Experience = req.Experience
	}
	if req.Education != "" {
		v.Education = req.Education
	}
	if req.Location != "" {
		v.Location = req.Location
	}
	if req.IsActive != nil {
		v.IsActive = *req.IsActive
	}
	v.UpdatedAt = time.Now()
}

func (v *Vacancy) Activate()   { b := true; v.IsActive = b; v.UpdatedAt = time.Now() }
func (v *Vacancy) Deactivate() { b := false; v.IsActive = b; v.UpdatedAt = time.Now() }
