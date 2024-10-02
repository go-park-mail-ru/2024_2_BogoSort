package storage

import (
	"errors"
	"fmt"
	"sync"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/services"
)

type Advert struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	ImageURL string `json:"image_url"`
	Price    uint   `json:"price"`
	Location string `json:"location"`
}

type AdvertsList struct {
	adverts  []*Advert
	advCount uint
	mu       sync.Mutex
}

var ErrAdvertNotFound = errors.New("объявление не найдено")

func (l *AdvertsList) Add(a *Advert) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.advCount++

	a.ID = l.advCount
	l.adverts = append(l.adverts, a)
}

func (l *AdvertsList) Update(a *Advert) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	for i, adv := range l.adverts {
		if adv.ID == a.ID {
			l.adverts[i] = a

			return nil
		}
	}

	return ErrAdvertNotFound
}

func (l *AdvertsList) DeleteAdvert(id uint) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	for i, adv := range l.adverts {
		if adv.ID == id {
			l.adverts = append(l.adverts[:i], l.adverts[i+1:]...)

			return nil
		}
	}

	return ErrAdvertNotFound
}

func (l *AdvertsList) GetAdverts() []Advert {
	l.mu.Lock()
	defer l.mu.Unlock()

	result := make([]Advert, len(l.adverts))

	for i, advert := range l.adverts {
		result[i] = *advert
	}

	return result
}

func (l *AdvertsList) GetAdvertByID(id uint) (Advert, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, advert := range l.adverts {
		if advert.ID == id {
			return *advert, nil
		}
	}

	return Advert{}, ErrAdvertNotFound
}

func NewAdvertsList() *AdvertsList {
	return &AdvertsList{
		adverts: make([]*Advert, 0),
		mu:      sync.Mutex{},
	}
}

func FillAdverts(ads *AdvertsList, imageService *services.ImageService) {
	ads.mu.Lock()
	defer ads.mu.Unlock()

	locations := []string{"Москва", "Санкт-Петербург", "Новосибирск", "Екатеринбург", "Казань"}
	titles := []string{
		"Продам велосипед", "Аренда квартиры", "Продам ноутбук", "Продам автомобиль", "Продам мебель",
		"Продам телефон", "Продам дом", "Аренда гаража", "Продам планшет", "Продам телевизор",
		"Продам велосипед", "Аренда квартиры", "Продам ноутбук", "Продам автомобиль", "Продам мебель",
		"Продам телефон", "Продам дом", "Аренда гаража", "Продам планшет", "Продам телевизор",
		"Продам велосипед", "Аренда квартиры", "Продам ноутбук", "Продам автомобиль", "Продам мебель",
		"Продам телефон", "Продам дом", "Аренда гаража", "Продам планшет", "Продам телевизор",
		"Продам велосипед", "Аренда квартиры", "Продам ноутбук", "Продам автомобиль", "Продам мебель",
		"Продам телефон", "Продам дом", "Аренда гаража", "Продам планшет", "Продам телевизор",
		"Продам велосипед", "Аренда квартиры", "Продам ноутбук", "Продам автомобиль", "Продам мебель",
		"Продам телефон", "Продам дом", "Аренда гаража", "Продам планшет", "Продам телевизор",
		"Продам велосипед", "Аренда квартиры", "Продам ноутбук", "Продам автомобиль", "Продам мебель",
		"Продам телефон", "Продам дом", "Аренда гаража", "Продам планшет", "Продам телевизор",
		"Продам велосипед", "Аренда квартиры", "Продам ноутбук", "Продам автомобиль", "Продам мебель",
		"Продам телефон", "Продам дом", "Аренда гаража", "Продам планшет", "Продам телевизор",
	}

	const testAdvCount, testPrice = 30, 1000

	for i := 1; i <= testAdvCount; i++ {
		imageURL := fmt.Sprintf("/static/images/image%d.jpg", i)

		id := uint(i)
		price := uint(testPrice + (i-1)*testPrice/10)

		advert := &Advert{
			ID:       id,
			Title:    titles[(i-1)%len(titles)],
			ImageURL: imageURL,
			Price:    price,
			Location: locations[(i-1)%len(locations)],
		}

		ads.adverts = append(ads.adverts, advert)

		imageService.SetImageURL(id, imageURL)
	}

	ads.advCount = uint(len(ads.adverts))
}
