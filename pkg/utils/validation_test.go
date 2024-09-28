package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	validEmails := []string{
		"test@example.com",
		"user.name+tag+sorting@example.com",
		"x@example.com",
		"example-indeed@strange-example.com",
	}

	invalidEmails := []string{
		"plainaddress",
		"@missingusername.com",
		"username@.com",
		"username@.com.",
	}

	for _, email := range validEmails {
		assert.True(t, ValidateEmail(email), "Email should be valid: %s", email)
	}

	for _, email := range invalidEmails {
		assert.False(t, ValidateEmail(email), "Email should be invalid: %s", email)
	}
}

func TestValidatePassword(t *testing.T) {
	validPasswords := []string{
		"Valid1Password!",
		"Another$Valid2",
	}

	invalidPasswords := []string{
		"short1!",
		"nouppercase1!",
		"NOLOWERCASE1!",
		"NoNumber!",
		"NoSpecialChar1",
	}

	for _, password := range validPasswords {
		assert.NoError(t, ValidatePassword(password), "Password should be valid: %s", password)
	}

	for _, password := range invalidPasswords {
		assert.Error(t, ValidatePassword(password), "Password should be invalid: %s", password)
	}
}