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
	list := storage.NewAdvertsList()
	imageService := services.NewImageService()
	handler := &AdvertsHandler{List: list, ImageService: imageService}

	req, err := http.NewRequest("GET", "/adverts", nil)
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
	list := storage.NewAdvertsList()
	imageService := services.NewImageService()
	handler := &AdvertsHandler{List: list, ImageService: imageService}

	advert := &storage.Advert{ID: 1, Title: "Test Advert"}
	list.Add(advert)

	req, err := http.NewRequest("GET", "/adverts/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/adverts/{id}", handler.GetAdvertByIDHandler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestAdvertsHandler_AddAdvertHandler(t *testing.T) {
	list := storage.NewAdvertsList()
	imageService := services.NewImageService()
	handler := &AdvertsHandler{List: list, ImageService: imageService}

	advertJSON := `{"title": "New Advert", "price": 1000, "location": "Москва"}`
	req, err := http.NewRequest("POST", "/adverts", strings.NewReader(advertJSON))
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
	list := storage.NewAdvertsList()
	imageService := services.NewImageService()
	handler := &AdvertsHandler{List: list, ImageService: imageService}

	advert := &storage.Advert{ID: 1, Title: "Test Advert"}
	list.Add(advert)

	updatedAdvertJSON := `{"id": 1, "title": "Updated Advert", "price": 2000, "location": "Санкт-Петербург"}`
	req, err := http.NewRequest("PUT", "/adverts/1", strings.NewReader(updatedAdvertJSON))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/adverts/{id}", handler.UpdateAdvertHandler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestAdvertsHandler_DeleteAdvertHandler(t *testing.T) {
	list := storage.NewAdvertsList()
	imageService := services.NewImageService()
	handler := &AdvertsHandler{List: list, ImageService: imageService}

	advert := &storage.Advert{ID: 1, Title: "Test Advert"}
	list.Add(advert)

	req, err := http.NewRequest("DELETE", "/adverts/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/adverts/{id}", handler.DeleteAdvertHandler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}
