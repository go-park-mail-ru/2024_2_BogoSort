package storage

import (
	"testing"
)

func TestAddSession(t *testing.T) {
	storage := NewSessionStorage()

	storage.AddSession("user@example.com", "token")
	if !storage.SessionExists("user@example.com") {
		t.Errorf("expected session to exist for user@example.com")
	}
}

func TestRemoveSession(t *testing.T) {
	storage := NewSessionStorage()

	storage.AddSession("user@example.com", "token")
	storage.RemoveSession("user@example.com")
	if storage.SessionExists("user@example.com") {
		t.Errorf("expected session to be removed for user@example.com")
	}
}

func TestSessionExists(t *testing.T) {
	storage := NewSessionStorage()

	storage.AddSession("user@example.com", "token")
	if !storage.SessionExists("user@example.com") {
		t.Errorf("expected session to exist for user@example.com")
	}
}
