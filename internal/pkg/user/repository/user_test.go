package repository

// import (
// 	"testing"
// )

// const (
// 	testEmail    = "newuser@example.com"
// 	testPassword = "passworD!1"
// )

// func TestCreateUser(t *testing.T) {
// 	storage := NewUserStorage()

// 	user, err := storage.CreateUser(testEmail, testPassword)
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}

// 	if user.Email != testEmail {
// 		t.Errorf("expected email to be %v, got %v", testEmail, user.Email)
// 	}
// }

// func TestGetUserByEmail(t *testing.T) {
// 	storage := NewUserStorage()

// 	_, err := storage.CreateUser(testEmail, testPassword)
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}

// 	user, err := storage.GetUserByEmail(testEmail)
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}

// 	if user.Email != testEmail {
// 		t.Errorf("expected email to be %v, got %v", testEmail, user.Email)
// 	}
// }

// func TestValidateUserByEmailAndPassword(t *testing.T) {
// 	storage := NewUserStorage()

// 	_, err := storage.CreateUser(testEmail, testPassword)
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}

// 	user, err := storage.ValidateUserByEmailAndPassword(testEmail, testPassword)
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}

// 	if user.Email != testEmail {
// 		t.Errorf("expected email to be %v, got %v", testEmail, user.Email)
// 	}
// }
