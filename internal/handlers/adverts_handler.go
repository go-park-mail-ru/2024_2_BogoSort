package myhandlers

import (
	"emporium/internal/models"
	"encoding/json"
	"net/http"
	"sync"
)

var (
	adverts         = []models.Advert{}
	advertIDCounter = 0
	mu              sync.Mutex
)

func GetadvertsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Advert{})
}

func AddTestAdvert() {
	mu.Lock()
	defer mu.Unlock()

	advert := models.Advert{
		ID:      1,
		Title:   "Test advert",
		Content: "This is a test advert.",
	}
	adverts = append(adverts, advert)
	advertIDCounter++
}
