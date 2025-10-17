package domain

import "errors"

var (
	ErrForbidden         = errors.New("forbidden")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrInvariantViolated = errors.New("invariant violated")
)
