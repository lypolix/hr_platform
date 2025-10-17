package bcrypt

import "golang.org/x/crypto/bcrypt"

type bcryptPasswordService struct {
	cost int
}

func NewBcryptPasswordService(cost int) *bcryptPasswordService {
	return &bcryptPasswordService{cost}
}

func (s *bcryptPasswordService) Hash(raw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(raw), s.cost)
	return string(bytes), err
}

func (s *bcryptPasswordService) Check(raw, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw))
	return err == nil
}
