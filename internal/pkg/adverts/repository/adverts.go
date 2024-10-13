package repository

import (
	"errors"
	"fmt"
	"sync"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/pkg/adverts/models"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/pkg/services"
)

var ErrAdvertNotFound = errors.New("объявление не найдено")

func (l *models.AdvertsList) CreateAdvert(a *models.Advert) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.advCount++

	a.ID = l.advCount
	l.adverts = append(l.adverts, a)
}

func (l *models.AdvertsList) UpdateAdvert(a *models.Advert) error {
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

func (l *models.AdvertsList) DeleteAdvert(id uint) error {
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

func (l *models.AdvertsList) GetAllAdverts() ([]models.Advert, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	result := make([]models.Advert, len(l.adverts))

	for i, advert := range l.adverts {
		result[i] = *advert
	}

	return result, nil
}

func (l *models.AdvertsList) GetAdvertById(id uint) (models.Advert, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, advert := range l.adverts {
		if advert.ID == id {
			return *advert, nil
		}
	}

	return models.Advert{}, ErrAdvertNotFound
}

func NewAdvertsList() *models.AdvertsList {
	return &models.AdvertsList{
		adverts: make([]*models.Advert, 0),
		mu:      sync.Mutex{},
	}
}

func FillAdverts(ads *models.AdvertsList, imageService *services.ImageService) {
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

		advert := &models.Advert{
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
