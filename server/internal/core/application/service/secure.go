package service

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordService struct{}

func NewBcryptPasswordService() *BcryptPasswordService { return &BcryptPasswordService{} }

func (BcryptPasswordService) Hash(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(b), err
}

func (BcryptPasswordService) Compare(hash, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}
