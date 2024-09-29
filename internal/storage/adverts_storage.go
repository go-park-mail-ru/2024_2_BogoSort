package storage

import (
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
	return fmt.Errorf("объявление не найдено")
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
	return fmt.Errorf("объявление не найдено")
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
	return Advert{}, fmt.Errorf("объявление не найдено")
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

	const testAdvCount = 60
	for i := 1; i <= testAdvCount; i++ {
		imageURL := fmt.Sprintf("/static/images/image%d.jpg", i)
		advert := &Advert{
			ID:       uint(i),
			Title:    titles[i%len(titles)],
			ImageURL: imageURL,
			Price:    uint(1000 + i*100),
			Location: locations[i%len(locations)],
		}
		ads.adverts = append(ads.adverts, advert)

		imageService.SetImageURL(uint(i)%10+1, imageURL)
	}

	ads.advCount = uint(len(ads.adverts))
}
