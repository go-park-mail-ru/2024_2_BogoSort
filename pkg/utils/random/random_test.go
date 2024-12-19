package random

import (
	"testing"
)

func TestBytesWithCharset(t *testing.T) {
	length := 10
	result, err := bytesWithCharset(length, charset)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != length {
		t.Fatalf("expected length %d, got %d", length, len(result))
	}
	for _, b := range result {
		if !contains(charset, b) {
			t.Fatalf("unexpected byte %c in result", b)
		}
	}
}

func TestBytes(t *testing.T) {
	length := 10
	result, err := Bytes(length)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != length {
		t.Fatalf("expected length %d, got %d", length, len(result))
	}
	for _, b := range result {
		if !contains(charset, b) {
			t.Fatalf("unexpected byte %c in result", b)
		}
	}
}

func contains(charset string, b byte) bool {
	for i := 0; i < len(charset); i++ {
		if charset[i] == b {
			return true
		}
	}
	return false
}
