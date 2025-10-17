package port

import "github.com/google/uuid"

type PasswordService interface {
	Hash(string) (string, error)
	Check(raw, hash string) bool
	Compare(hash, plain string) error
}

type role string

const (
	RoleCompany    role = "company"
	RoleUniversity role = "university"
	RoleAdmin      role = "admin"
)

type TokenPayload struct {
	Sub  uuid.UUID
	Role role
}

type TokenService interface {
	Generate(TokenPayload) (string, error)
	Validate(string) (TokenPayload, error)
}

