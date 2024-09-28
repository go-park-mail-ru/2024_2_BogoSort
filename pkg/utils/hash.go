package utils

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 14
)

func HashPassword(password string) (string) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)

	if err != nil {
		return ""
	}

	return string(bytes)
}

func ComparePassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
