package services

import (
	"testing"
)

func TestImageService_GetImageURL(t *testing.T) {
	service := NewImageService()
	service.SetImageURL(1, "/static/images/image1.jpg")

	url, err := service.GetImageURL(1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if url != "/static/images/image1.jpg" {
		t.Fatalf("expected '/static/images/image1.jpg', got %v", url)
	}

	_, err = service.GetImageURL(2)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestImageService_SetImageURL(t *testing.T) {
	service := NewImageService()
	service.SetImageURL(1, "/static/images/image1.jpg")

	url, err := service.GetImageURL(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	
	if url != "/static/images/image1.jpg" {
		t.Fatalf("expected '/static/images/image1.jpg', got %v", url)
	}
}
