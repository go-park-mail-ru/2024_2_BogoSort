package storage

import (
	"fmt"
	"sync"
)

type Advert struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type AdvertsList struct {
	adverts []Advert
}

var (
	adv      = AdvertsList{}
	advCount uint
	mu       sync.Mutex
)

func init() {
	FillAdverts()
}

func AddAdvert(advert Advert) {
	mu.Lock()
	defer mu.Unlock()
	adv.adverts = append(adv.adverts, advert)
	advCount++
}

func GetAdverts() []Advert {
	mu.Lock()
	defer mu.Unlock()
	return adv.adverts
}

func GetAdvertByID(id uint) (Advert, error) {
	mu.Lock()
	defer mu.Unlock()
	for _, advert := range adv.adverts {
		if advert.ID == id {
			return advert, nil
		}
	}
	return Advert{}, fmt.Errorf("advert not found")
}

func FillAdverts() {
	AddAdvert(Advert{ID: 1, Title: "First advert", Content: "This is the first advert"})
	AddAdvert(Advert{ID: 2, Title: "Second advert", Content: "This is the second advert"})
	AddAdvert(Advert{ID: 3, Title: "Third advert", Content: "This is the third advert"})
	advCount = 3
}
