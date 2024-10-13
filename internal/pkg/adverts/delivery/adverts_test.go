package delivery

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/services"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/storage"
	"github.com/stretchr/testify/assert"

	"github.com/gorilla/mux"
)

func TestAdvertsHandler_GetAdvertsHandler(t *testing.T) {
	t.Parallel()

	list := storage.NewAdvertsList()
	imageService := services.NewImageService()
	handler := &AdvertsHandler{List: list, ImageService: imageService}

	req, err := http.NewRequest(http.MethodGet, "/api/v1/adverts", nil)
	assert.NoError(t, err, "failed to create request")

	rr := httptest.NewRecorder()
	handler.GetAdvertsHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
}

func TestAdvertsHandler_GetAdvertByIDHandler(t *testing.T) {
	t.Parallel()

	list := storage.NewAdvertsList()
	imageService := services.NewImageService()
	handler := &AdvertsHandler{List: list, ImageService: imageService}

	advert := &storage.Advert{ID: 1, Title: "Test Advert"}
	list.Add(advert)

	t.Run("Valid ID", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/v1/adverts/1", nil)
		assert.NoError(t, err, "failed to create request")

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/api/v1/adverts/{id}", handler.GetAdvertByIDHandler)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	})

	t.Run("Invalid ID", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/v1/adverts/invalid", nil)
		assert.NoError(t, err, "failed to create request")

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/api/v1/adverts/{id}", handler.GetAdvertByIDHandler)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")
		assert.Contains(t, rr.Body.String(), ErrInvalidID.Error(), "expected invalid ID error")
	})

	t.Run("Non-existent ID", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/v1/adverts/999", nil)
		assert.NoError(t, err, "failed to create request")

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/api/v1/adverts/{id}", handler.GetAdvertByIDHandler)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code, "handler returned wrong status code")
		assert.Contains(t, rr.Body.String(), ErrAdvertNotFound.Error(), "expected advert not found error")
	})
}

func TestAdvertsHandler_AddAdvertHandler(t *testing.T) {
	t.Parallel()

	list := storage.NewAdvertsList()
	imageService := services.NewImageService()
	handler := &AdvertsHandler{List: list, ImageService: imageService}

	t.Run("Valid Advert", func(t *testing.T) {
		advertJSON := `{"title": "New Advert", "price": 1000, "location": "Москва"}`
		req, err := http.NewRequest(http.MethodPost, "/api/v1/adverts", strings.NewReader(advertJSON))
		assert.NoError(t, err, "failed to create request")

		rr := httptest.NewRecorder()
		handler.AddAdvertHandler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/api/v1/adverts", strings.NewReader("invalid json"))
		assert.NoError(t, err, "failed to create request")

		rr := httptest.NewRecorder()
		handler.AddAdvertHandler(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")
		assert.Contains(t, rr.Body.String(), ErrFailedToAddAdvert.Error(), "expected failed to add advert error")
	})
}

func TestAdvertsHandler_UpdateAdvertHandler(t *testing.T) {
	t.Parallel()

	list := storage.NewAdvertsList()
	imageService := services.NewImageService()
	handler := &AdvertsHandler{List: list, ImageService: imageService}

	advert := &storage.Advert{ID: 1, Title: "Test Advert"}
	list.Add(advert)

	t.Run("Valid Update", func(t *testing.T) {
		updatedAdvertJSON := `{"id": 1, "title": "Updated Advert", "price": 2000, "location": "Санкт-Петербург"}`
		req, err := http.NewRequest(http.MethodPut, "/api/v1/adverts/1", strings.NewReader(updatedAdvertJSON))
		assert.NoError(t, err, "failed to create request")

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/api/v1/adverts/{id}", handler.UpdateAdvertHandler)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	})

	t.Run("Invalid ID", func(t *testing.T) {
		updatedAdvertJSON := `{"id": 1, "title": "Updated Advert", "price": 2000, "location": "Санкт-Петербург"}`
		req, err := http.NewRequest(http.MethodPut, "/api/v1/adverts/invalid", strings.NewReader(updatedAdvertJSON))
		assert.NoError(t, err, "failed to create request")

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/api/v1/adverts/{id}", handler.UpdateAdvertHandler)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")
		assert.Contains(t, rr.Body.String(), ErrInvalidID.Error(), "expected invalid ID error")
	})

	t.Run("Mismatched ID", func(t *testing.T) {
		updatedAdvertJSON := `{"id": 2, "title": "Updated Advert", "price": 2000, "location": "Санкт-Петербург"}`
		req, err := http.NewRequest(http.MethodPut, "/api/v1/adverts/1", strings.NewReader(updatedAdvertJSON))
		assert.NoError(t, err, "failed to create request")

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/api/v1/adverts/{id}", handler.UpdateAdvertHandler)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")
		assert.Contains(t, rr.Body.String(), "Id in URL and JSON do not match", "expected ID mismatch error")
	})

	t.Run("Non-existent ID", func(t *testing.T) {
		updatedAdvertJSON := `{"id": 999, "title": "Updated Advert", "price": 2000, "location": "Санкт-Петербург"}`
		req, err := http.NewRequest(http.MethodPut, "/api/v1/adverts/999", strings.NewReader(updatedAdvertJSON))
		assert.NoError(t, err, "failed to create request")

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/api/v1/adverts/{id}", handler.UpdateAdvertHandler)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code, "handler returned wrong status code")
		assert.Contains(t, rr.Body.String(), ErrFailedToUpdateAdvert.Error(), "expected failed to update advert error")
	})
}

func TestAdvertsHandler_DeleteAdvertHandler(t *testing.T) {
	t.Parallel()

	list := storage.NewAdvertsList()
	imageService := services.NewImageService()
	handler := &AdvertsHandler{List: list, ImageService: imageService}

	advert := &storage.Advert{ID: 1, Title: "Test Advert"}
	list.Add(advert)

	t.Run("Valid Delete", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, "/api/v1/adverts/1", nil)
		assert.NoError(t, err, "failed to create request")

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/api/v1/adverts/{id}", handler.DeleteAdvertHandler)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code, "handler returned wrong status code")
	})

	t.Run("Invalid ID", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, "/api/v1/adverts/invalid", nil)
		assert.NoError(t, err, "failed to create request")

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/api/v1/adverts/{id}", handler.DeleteAdvertHandler)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")
		assert.Contains(t, rr.Body.String(), ErrInvalidID.Error(), "expected invalid ID error")
	})

	t.Run("Non-existent ID", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, "/api/v1/adverts/999", nil)
		assert.NoError(t, err, "failed to create request")

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/api/v1/adverts/{id}", handler.DeleteAdvertHandler)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code, "handler returned wrong status code")
		assert.Contains(t, rr.Body.String(), ErrFailedToDeleteAdvert.Error(), "expected failed to delete advert error")
	})
}
