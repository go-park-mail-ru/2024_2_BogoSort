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
	adverts []*Advert
	mu      sync.Mutex
}

var (
	advertsList = &AdvertsList{}
	advCount    uint
)

func (l *AdvertsList) Add(a *Advert) {
	l.mu.Lock()
	defer l.mu.Unlock()
	advCount++
	a.ID = advCount
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

func FillAdverts(ads *AdvertsList) {
	ads.mu.Lock()
	defer ads.mu.Unlock()

	testAdverts := []*Advert{
		{
			ID:      1,
			Title:   "Продается автомобиль",
			Content: "Хорошее состояние, небольшой пробег",
		},
		{
			ID:      2,
			Title:   "Сдается квартира",
			Content: "2 комнаты, центр города",
		},
		{
			ID:      3,
			Title:   "Продам ноутбук",
			Content: "Почти новый, игровой",
		},
	}

	ads.adverts = append(ads.adverts, testAdverts...)
}
