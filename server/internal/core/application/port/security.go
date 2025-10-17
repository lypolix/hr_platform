package port

import "github.com/google/uuid"

type PasswordService interface {
	Hash(string) (string, error)
	Check(raw, hash string) bool
	Compare(hash, plain string) error
}

type Role string

const (
	RoleCompany    Role = "company"
	RoleUniversity Role = "university"
	RoleAdmin      Role = "admin"
)

type TokenPayload struct {
	Sub  uuid.UUID
	Role Role
}

type TokenService interface {
	Generate(TokenPayload) (string, error)
	Validate(string) (TokenPayload, error)
}
