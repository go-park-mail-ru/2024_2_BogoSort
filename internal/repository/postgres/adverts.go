package postgres

import (
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/services"
	"sync"
)

var ErrAdvertNotFound = errors.New("advert not found")

type advertRepository struct {
	adverts  []*entity.Advert
	advCount uint
	mu       sync.Mutex
}

func NewAdvertRepository() entity.AdvertRepository {
	return &advertRepository{
		adverts:  make([]*entity.Advert, 0),
		advCount: 0,
		mu:       sync.Mutex{},
	}
}

func (l *advertRepository) CreateAdvert(a *entity.Advert) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.advCount++

	a.ID = l.advCount
	l.adverts = append(l.adverts, a)
	return nil
}

func (l *advertRepository) UpdateAdvert(a *entity.Advert) error {
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

func (l *advertRepository) DeleteAdvert(id uint) error {
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

func (l *advertRepository) GetAllAdverts() ([]*entity.Advert, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	result := make([]*entity.Advert, len(l.adverts))
	copy(result, l.adverts)

	return result, nil
}

func (l *advertRepository) GetAdvertById(id uint) (*entity.Advert, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, advert := range l.adverts {
		if advert.ID == id {
			return advert, nil
		}
	}

	return nil, ErrAdvertNotFound
}

func (l *advertRepository) NewAdvertsList() *entity.AdvertsList {
	return &entity.AdvertsList{
		Adverts:  make([]*entity.Advert, 0),
		AdvCount: 0,
		Mu:       sync.Mutex{},
	}
}

func (l *advertRepository) FillAdverts(ads *entity.AdvertsList, imageService *services.ImageService) {
	ads.Mu.Lock()
	defer ads.Mu.Unlock()

	locations := []string{"Москва", "Санкт-Петербург", "Новосибирск", "Екатеринбург", "Казань"}
	titles := []string{
		"Продам велосипед", "Аренда квартиры", "Продам ноутбук", "Продам автомобиль", "Продам мебель",
		"Продам телефон", "Продам дом", "Аренда гаража", "Продам планшет", "Продам телевизор",
	}

	const testAdvCount, testPrice = 30, 1000

	for i := 1; i <= testAdvCount; i++ {
		imageURL := fmt.Sprintf("/static/images/image%d.jpg", i)
		id := uint(i)
		price := uint(testPrice + (i-1)*testPrice/10)

		advert := &entity.Advert{
			ID:       id,
			Title:    titles[(i-1)%len(titles)],
			ImageURL: imageURL,
			Price:    price,
			Location: locations[(i-1)%len(locations)],
		}

		ads.Adverts = append(ads.Adverts, advert)
		imageService.SetImageURL(id, imageURL)
	}

	ads.AdvCount = uint(len(ads.Adverts))
}
