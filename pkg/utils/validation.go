package utils

import (
	"errors"
	"unicode"
)

const (
	minPasswordLength = 8
	maxPasswordLength = 32
)

var (
	ErrPasswordTooShort     = errors.New("password is too short")
	ErrPasswordTooLong      = errors.New("password is too long")
	ErrPasswordRequirements = errors.New("password must contain upper, lower, number, and special characters")
)

func ValidatePassword(password string) error {
	if err := checkPasswordLength(password); err != nil {
		return err
	}

	if err := checkPasswordComplexity(password); err != nil {
		return err
	}

	return nil
}

func checkPasswordLength(password string) error {
	if len(password) < minPasswordLength {
		return ErrPasswordTooShort
	}

	if len(password) > maxPasswordLength {
		return ErrPasswordTooLong
	}

	return nil
}

func checkPasswordComplexity(password string) error {
	var hasUpper, hasLower, hasNumber, hasSpecial bool

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
