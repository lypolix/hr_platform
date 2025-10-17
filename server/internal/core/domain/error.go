package domain

import "errors"

var (
	ErrConflict          = errors.New("conflict")
	ErrNotFound          = errors.New("not found")
	ErrForbidden         = errors.New("forbidden")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrInvariantViolated = errors.New("invariant violated")
)
