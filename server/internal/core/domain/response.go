package domain

import (
	"time"

	"github.com/google/uuid"
)

// Статусы отклика
const (
	ResponseStatusNew      = "new"
	ResponseStatusViewed   = "viewed"
	ResponseStatusRejected = "rejected"
	ResponseStatusAccepted = "accepted"
)

type Response struct {
	ID        uuid.UUID `json:"id" db:"id"`
	VacancyID uuid.UUID `json:"vacancy_id" db:"vacancy_id"`

	FullName string `json:"full_name" db:"full_name"`
	Email    string `json:"email" db:"email"`
	Phone    string `json:"phone" db:"phone"`

	CoverLetter string `json:"cover_letter" db:"cover_letter"`
	ResumeURL   string `json:"resume_url" db:"resume_url"` // ссылка на S3-объект

	Status    string    `json:"status" db:"status"` // new/viewed/rejected/accepted
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateResponseRequest struct {
	VacancyID   uuid.UUID `json:"vacancy_id" binding:"required"`
	FullName    string    `json:"full_name" binding:"required"`
	Email       string    `json:"email" binding:"required,email"`
	Phone       string    `json:"phone" binding:"required"`
	CoverLetter string    `json:"cover_letter"`
	ResumeURL   string    `json:"resume_url" binding:"required"` // путь/URL, файл уже загружен в S3
}

type UpdateResponseStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=new viewed rejected accepted"`
}

type ResponseFilter struct {
	VacancyID *uuid.UUID
	Status    *string
	Limit     int
	Offset    int
}

func NewResponse(req CreateResponseRequest) *Response {
	now := time.Now()
	return &Response{
		ID:          uuid.New(),
		VacancyID:   req.VacancyID,
		FullName:    req.FullName,
		Email:       req.Email,
		Phone:       req.Phone,
		CoverLetter: req.CoverLetter,
		ResumeURL:   req.ResumeURL,
		Status:      ResponseStatusNew,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (r *Response) UpdateStatus(status string) {
	r.Status = status
	r.UpdatedAt = time.Now()
}
