package util

import (
	"golang.org/x/crypto/bcrypt"
)

type Bcrypt interface {
	GenerateFromPassword([]byte, int) ([]byte, error)
	CompareHashAndPassword([]byte, []byte) error
	GetDefaultCost() int
}

type StandardBcrypt struct {
}

func (bc *StandardBcrypt) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

func (bc *StandardBcrypt) CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

func (bc *StandardBcrypt) GetDefaultCost() int {
	return bcrypt.DefaultCost
}
