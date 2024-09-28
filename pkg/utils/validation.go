package utils

import (
	"errors"
	"regexp"
	"unicode"
)

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

const (
	minPasswordLength = 8
	maxPasswordLength = 64
)

var (
	ErrPasswordTooShort     = errors.New("password is too short")
	ErrPasswordTooLong      = errors.New("password is too long")
	ErrPasswordRequirements = errors.New("password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
)

func ValidatePassword(password string) error {
	if len(password) < minPasswordLength {
		return ErrPasswordTooShort
	}
	if len(password) > maxPasswordLength {
		return ErrPasswordTooLong
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return ErrPasswordRequirements
	}

	return nil
}

