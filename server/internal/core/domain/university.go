package domain

import "github.com/google/uuid"

type (
	University struct {
		id           uuid.UUID
		title        string
		login        string
		passwordHash string
		inn          string
		confirmed    bool
	}
)
