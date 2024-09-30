package utils

import (
	"os"
	"testing"

	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET_KEY", "test_secret_key")
	os.Setenv("JWT_EXPIRATION_TIME", "1h")
	os.Setenv("JWT_ISSUER", "test_issuer")

	config.InitFromEnv()
	InitJWT()

	code := m.Run()
	os.Exit(code)
}

func TestCreateToken(t *testing.T) {
	email := "test@example.com"
	token, err := CreateToken(email)
	assert.NoError(t, err, "Token creation should not return an error")
	assert.NotEmpty(t, token, "Token should not be empty")
}

func TestValidateToken(t *testing.T) {
	email := "test@example.com"
	token, err := CreateToken(email)
	assert.NoError(t, err, "Token creation should not return an error")

	subject, err := ValidateToken(token)
	assert.NoError(t, err, "Token validation should not return an error")
	assert.Equal(t, email, subject, "Token subject should match the email")
}

func TestValidateToken_Invalid(t *testing.T) {
	_, err := ValidateToken("invalid_token")
	assert.Error(t, err, "Invalid token should return an error")
}
