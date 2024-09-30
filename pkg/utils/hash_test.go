package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "TestPassword123!"
	hash := HashPassword(password)
	assert.NotEmpty(t, hash, "Hash should not be empty")
}

func TestComparePassword(t *testing.T) {
	password := "TestPassword123!"
	hash := HashPassword(password)
	assert.True(t, ComparePassword(password, hash), "Password should match the hash")
	assert.False(t, ComparePassword("WrongPassword", hash), "Password should not match the hash")
}
