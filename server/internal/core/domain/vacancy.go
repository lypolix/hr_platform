package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type (
	Vacancy struct {
		id          uuid.UUID
		companyID   uuid.UUID
		title       string
		description string
		contacts    string

		requirements     string
		responsibilities string
		conditions       string

		salaryFrom *int
		salaryTo   *int

		employment string
		schedule   string
		experience string
		education  string
		location   string

		isActive  bool
		createdAt time.Time
		updatedAt time.Time
	}

	VacancyImmutable struct {
		ID               uuid.UUID
		CompanyID        uuid.UUID
		Title            string
		Description      string
		Contacts         string
		Requirements     string
		Responsibilities string
		Conditions       string
		SalaryFrom       *int
		SalaryTo         *int
		Employment       string
		Schedule         string
		Experience       string
		Education        string
		Location         string
		IsActive         bool
		CreatedAt        time.Time
		UpdatedAt        time.Time
	}

	CreateVacancyAttrs struct {
		CompanyID       uuid.UUID
		Title           string
		Description     string
		Contacts        string
		Requirements     string
		Responsibilities string
		Conditions       string
		SalaryFrom       *int
		SalaryTo         *int
		Employment       string
		Schedule         string
		Experience       string
		Education        string
		Location         string
	}

	VacancyFilter struct {
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
)

func (v *Vacancy) Immutable() VacancyImmutable {
	return VacancyImmutable{
		ID:               v.id,
		CompanyID:        v.companyID,
		Title:            v.title,
		Description:      v.description,
		Contacts:         v.contacts,
		Requirements:     v.requirements,
		Responsibilities: v.responsibilities,
		Conditions:       v.conditions,
		SalaryFrom:       v.salaryFrom,
		SalaryTo:         v.salaryTo,
		Employment:       v.employment,
		Schedule:         v.schedule,
		Experience:       v.experience,
		Education:        v.education,
		Location:         v.location,
		IsActive:         v.isActive,
		CreatedAt:        v.createdAt,
		UpdatedAt:        v.updatedAt,
	}
}

func (v *Vacancy) checkInvariants() error {
	if v.id == uuid.Nil {
		return fmt.Errorf("%w: nil id", ErrInvariantViolated)
	}
	if v.companyID == uuid.Nil {
		return fmt.Errorf("%w: nil company id", ErrInvariantViolated)
	}
	if l := len(v.title); l < 1 || l > 512 {
		return fmt.Errorf("%w: invalid title length", ErrInvariantViolated)
	}
	if l := len(v.description); l < 1 {
		return fmt.Errorf("%w: empty description", ErrInvariantViolated)
	}
	if v.salaryFrom != nil && *v.salaryFrom < 0 {
		return fmt.Errorf("%w: negative salary_from", ErrInvariantViolated)
	}
	if v.salaryTo != nil && *v.salaryTo < 0 {
		return fmt.Errorf("%w: negative salary_to", ErrInvariantViolated)
	}
	if v.salaryFrom != nil && v.salaryTo != nil && *v.salaryFrom > *v.salaryTo {
		return fmt.Errorf("%w: salary_from > salary_to", ErrInvariantViolated)
	}
	if v.createdAt.IsZero() {
		return fmt.Errorf("%w: zero creation time", ErrInvariantViolated)
	}
	if v.updatedAt.IsZero() {
		return fmt.Errorf("%w: zero updation time", ErrInvariantViolated)
	}
	if v.createdAt.After(v.updatedAt) {
		return fmt.Errorf("%w: creation time is after updation time", ErrInvariantViolated)
	}
	return nil
}

func CreateVacancy(attrs CreateVacancyAttrs, at time.Time) (*Vacancy, error) {
	imm := VacancyImmutable{
		ID:               uuid.New(),
		CompanyID:        attrs.CompanyID,
		Title:            attrs.Title,
		Description:      attrs.Description,
		Contacts:         attrs.Contacts,
		Requirements:     attrs.Requirements,
		Responsibilities: attrs.Responsibilities,
		Conditions:       attrs.Conditions,
		SalaryFrom:       attrs.SalaryFrom,
		SalaryTo:         attrs.SalaryTo,
		Employment:       attrs.Employment,
		Schedule:         attrs.Schedule,
		Experience:       attrs.Experience,
		Education:        attrs.Education,
		Location:         attrs.Location,
		IsActive:         true,
		CreatedAt:        at,
		UpdatedAt:        at,
	}
	return ReconstructVacancy(imm)
}

func ReconstructVacancy(immutable VacancyImmutable) (*Vacancy, error) {
	v := &Vacancy{
		id:               immutable.ID,
		companyID:        immutable.CompanyID,
		title:            immutable.Title,
		description:      immutable.Description,
		contacts:         immutable.Contacts,
		requirements:     immutable.Requirements,
		responsibilities: immutable.Responsibilities,
		conditions:       immutable.Conditions,
		salaryFrom:       immutable.SalaryFrom,
		salaryTo:         immutable.SalaryTo,
		employment:       immutable.Employment,
		schedule:         immutable.Schedule,
		experience:       immutable.Experience,
		education:       immutable.Education,
		location:         immutable.Location,
		isActive:         immutable.IsActive,
		createdAt:        immutable.CreatedAt,
		updatedAt:        immutable.UpdatedAt,
	}
	return v, v.checkInvariants()
}

// Мутации через Immutable + Reconstruct
func (v *Vacancy) Update(patch CreateVacancyAttrs, at time.Time) (*Vacancy, error) {
	imm := v.Immutable()
	if patch.Title != "" {
		imm.Title = patch.Title
	}
	if patch.Description != "" {
		imm.Description = patch.Description
	}
	if patch.Contacts != "" {
		imm.Contacts = patch.Contacts
	}
	if patch.Requirements != "" {
		imm.Requirements = patch.Requirements
	}
	if patch.Responsibilities != "" {
		imm.Responsibilities = patch.Responsibilities
	}
	if patch.Conditions != "" {
		imm.Conditions = patch.Conditions
	}
	if patch.SalaryFrom != nil {
		imm.SalaryFrom = patch.SalaryFrom
	}
	if patch.SalaryTo != nil {
		imm.SalaryTo = patch.SalaryTo
	}
	if patch.Employment != "" {
		imm.Employment = patch.Employment
	}
	if patch.Schedule != "" {
		imm.Schedule = patch.Schedule
	}
	if patch.Experience != "" {
		imm.Experience = patch.Experience
	}
	if patch.Education != "" {
		imm.Education = patch.Education
	}
	if patch.Location != "" {
		imm.Location = patch.Location
	}
	imm.UpdatedAt = at
	return ReconstructVacancy(imm)
}

func (v *Vacancy) Activate(at time.Time) (*Vacancy, error) {
	imm := v.Immutable()
	imm.IsActive = true
	imm.UpdatedAt = at
	return ReconstructVacancy(imm)
}

func (v *Vacancy) Deactivate(at time.Time) (*Vacancy, error) {
	imm := v.Immutable()
	imm.IsActive = false
	imm.UpdatedAt = at
	return ReconstructVacancy(imm)
}
