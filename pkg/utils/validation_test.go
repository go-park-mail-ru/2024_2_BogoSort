package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{
			name:     "Valid password",
			password: "Valid1@Password",
			wantErr:  nil,
		},
		{
			name:     "Password too short",
			password: "Short1@",
			wantErr:  ErrPasswordTooShort,
		},
		{
			name:     "Password too long",
			password: "ThisIsAVeryLongPasswordThatExceedsTheMaximumLength1@",
			wantErr:  ErrPasswordTooLong,
		},
		{
			name:     "Missing uppercase letter",
			password: "valid1@password",
			wantErr:  ErrPasswordRequirements,
		},
		{
			name:     "Missing lowercase letter",
			password: "VALID1@PASSWORD",
			wantErr:  ErrPasswordRequirements,
		},
		{
			name:     "Missing number",
			password: "Valid@Password",
			wantErr:  ErrPasswordRequirements,
		},
		{
			name:     "Missing special character",
			password: "Valid1Password",
			wantErr:  ErrPasswordRequirements,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
