package storage

import (
	"testing"
)

func TestCreateUser(t *testing.T) {
	storage := NewUserStorage()

	user, err := storage.CreateUser("newuser@example.com", "password")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.Email != "newuser@example.com" {
		t.Errorf("expected email to be newuser@example.com, got %v", user.Email)
	}
}

func TestGetUserByEmail(t *testing.T) {
	storage := NewUserStorage()

	_, err := storage.CreateUser("newuser@example.com", "password")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	user, err := storage.GetUserByEmail("newuser@example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.Email != "newuser@example.com" {
		t.Errorf("expected email to be newuser@example.com, got %v", user.Email)
	}
}

func TestValidateUserByEmailAndPassword(t *testing.T) {
	storage := NewUserStorage()

	_, err := storage.CreateUser("newuser@example.com", "password")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	user, err := storage.ValidateUserByEmailAndPassword("newuser@example.com", "password")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.Email != "newuser@example.com" {
		t.Errorf("expected email to be newuser@example.com, got %v", user.Email)
	}
}
