package domain

import (
	"github.com/google/uuid"
)

type role string

const (
	RoleCompany    role = "company"
	RoleUniversity role = "university"
	RoleAdmin      role = "admin"
)

type Actor struct {
	ID   uuid.UUID
	Role role
}
