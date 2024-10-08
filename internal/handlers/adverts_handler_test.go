package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/services"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/storage"

	"github.com/gorilla/mux"
)

func TestAdvertsHandler_GetAdvertsHandler(t *testing.T) {
	t.Parallel()

	list := storage.NewAdvertsList()
	imageService := services.NewImageService()
	handler := &AdvertsHandler{List: list, ImageService: imageService}

	req, err := http.NewRequest(http.MethodGet, "/api/v1/adverts", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetAdvertsHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestAdvertsHandler_GetAdvertByIDHandler(t *testing.T) {
	t.Parallel()

	list := storage.NewAdvertsList()
	imageService := services.NewImageService()
	handler := &AdvertsHandler{List: list, ImageService: imageService}

	advert := &storage.Advert{ID: 1, Title: "Test Advert"}
	list.Add(advert)

	req, err := http.NewRequest(http.MethodGet, "/api/v1/adverts/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/adverts/{id}", handler.GetAdvertByIDHandler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestAdvertsHandler_AddAdvertHandler(t *testing.T) {
	t.Parallel()

	list := storage.NewAdvertsList()
	imageService := services.NewImageService()
	handler := &AdvertsHandler{List: list, ImageService: imageService}

	advertJSON := `{"title": "New Advert", "price": 1000, "location": "Москва"}`
	req, err := http.NewRequest(http.MethodPost, "/api/v1/adverts", strings.NewReader(advertJSON))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.AddAdvertHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestAdvertsHandler_UpdateAdvertHandler(t *testing.T) {
	t.Parallel()

	list := storage.NewAdvertsList()
	imageService := services.NewImageService()
	handler := &AdvertsHandler{List: list, ImageService: imageService}

	advert := &storage.Advert{ID: 1, Title: "Test Advert"}
	list.Add(advert)

	updatedAdvertJSON := `{"id": 1, "title": "Updated Advert", "price": 2000, "location": "Санкт-Петербург"}`

	req, err := http.NewRequest(http.MethodPut, "/api/v1/adverts/1", strings.NewReader(updatedAdvertJSON))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/adverts/{id}", handler.UpdateAdvertHandler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestAdvertsHandler_DeleteAdvertHandler(t *testing.T) {
	t.Parallel()

	list := storage.NewAdvertsList()
	imageService := services.NewImageService()
	handler := &AdvertsHandler{List: list, ImageService: imageService}

	advert := &storage.Advert{ID: 1, Title: "Test Advert"}
	list.Add(advert)

	req, err := http.NewRequest(http.MethodDelete, "/api/v1/adverts/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/adverts/{id}", handler.DeleteAdvertHandler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}
