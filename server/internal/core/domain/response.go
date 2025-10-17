package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	ResponseStatusNew      = "new"
	ResponseStatusViewed   = "viewed"
	ResponseStatusRejected = "rejected"
	ResponseStatusAccepted = "accepted"
)

type (
	Response struct {
		id        uuid.UUID
		vacancyID uuid.UUID
		fullName  string
		email     string
		phone     string
		coverLetter string
		resumeURL string
		status    string
		createdAt time.Time
		updatedAt time.Time
	}

	ResponseImmutable struct {
		ID          uuid.UUID
		VacancyID   uuid.UUID
		FullName    string
		Email       string
		Phone       string
		CoverLetter string
		ResumeURL   string
		Status      string
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	CreateResponseAttrs struct {
		VacancyID   uuid.UUID
		FullName    string
		Email       string
		Phone       string
		CoverLetter string
		ResumeURL   string
	}

	ResponseFilter struct {
		VacancyID *uuid.UUID
		Status    *string
		Limit     int
		Offset    int
	}
)

func (r *Response) Immutable() ResponseImmutable {
	return ResponseImmutable{
		ID:          r.id,
		VacancyID:   r.vacancyID,
		FullName:    r.fullName,
		Email:       r.email,
		Phone:       r.phone,
		CoverLetter: r.coverLetter,
		ResumeURL:   r.resumeURL,
		Status:      r.status,
		CreatedAt:   r.createdAt,
		UpdatedAt:   r.updatedAt,
	}
}

func (r *Response) checkInvariants() error {
	if r.id == uuid.Nil {
		return fmt.Errorf("%w: nil id", ErrInvariantViolated)
	}
	if r.vacancyID == uuid.Nil {
		return fmt.Errorf("%w: nil vacancy id", ErrInvariantViolated)
	}
	if l := len(r.fullName); l < 1 || l > 256 {
		return fmt.Errorf("%w: invalid full_name length", ErrInvariantViolated)
	}
	if l := len(r.email); l < 3 || l > 256 {
		return fmt.Errorf("%w: invalid email length", ErrInvariantViolated)
	}
	switch r.status {
	case ResponseStatusNew, ResponseStatusViewed, ResponseStatusRejected, ResponseStatusAccepted:
	default:
		return fmt.Errorf("%w: invalid status", ErrInvariantViolated)
	}
	if r.createdAt.IsZero() {
		return fmt.Errorf("%w: zero creation time", ErrInvariantViolated)
	}
	if r.updatedAt.IsZero() {
		return fmt.Errorf("%w: zero updation time", ErrInvariantViolated)
	}
	if r.createdAt.After(r.updatedAt) {
		return fmt.Errorf("%w: creation time is after updation time", ErrInvariantViolated)
	}
	return nil
}

func CreateResponse(attrs CreateResponseAttrs, at time.Time) (*Response, error) {
	imm := ResponseImmutable{
		ID:          uuid.New(),
		VacancyID:   attrs.VacancyID,
		FullName:    attrs.FullName,
		Email:       attrs.Email,
		Phone:       attrs.Phone,
		CoverLetter: attrs.CoverLetter,
		ResumeURL:   attrs.ResumeURL,
		Status:      ResponseStatusNew,
		CreatedAt:   at,
		UpdatedAt:   at,
	}
	return ReconstructResponse(imm)
}

func ReconstructResponse(immutable ResponseImmutable) (*Response, error) {
	r := &Response{
		id:          immutable.ID,
		vacancyID:   immutable.VacancyID,
		fullName:    immutable.FullName,
		email:       immutable.Email,
		phone:       immutable.Phone,
		coverLetter: immutable.CoverLetter,
		resumeURL:   immutable.ResumeURL,
		status:      immutable.Status,
		createdAt:   immutable.CreatedAt,
		updatedAt:   immutable.UpdatedAt,
	}
	return r, r.checkInvariants()
}

// Мутации через Immutable + Reconstruct
func (r *Response) SetStatus(status string, at time.Time) (*Response, error) {
	imm := r.Immutable()
	imm.Status = status
	imm.UpdatedAt = at
	return ReconstructResponse(imm)
}

func (r *Response) UpdateContacts(fullName, email, phone string, at time.Time) (*Response, error) {
	imm := r.Immutable()
	if fullName != "" {
		imm.FullName = fullName
	}
	if email != "" {
		imm.Email = email
	}
	if phone != "" {
		imm.Phone = phone
	}
	imm.UpdatedAt = at
	return ReconstructResponse(imm)
}
