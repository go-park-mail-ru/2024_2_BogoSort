package entity

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/services"
	"sync"
)

type Advert struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	ImageURL string `json:"image_url"`
	Price    uint   `json:"price"`
	Location string `json:"location"`
}

type AdvertsList struct {
	Adverts  []*Advert
	AdvCount uint
	Mu       sync.Mutex
}

type AdvertRepository interface {
	CreateAdvert(advert *Advert) error
	GetAllAdverts() ([]*Advert, error)
	GetAdvertById(id uint) (*Advert, error)
	UpdateAdvert(advert *Advert) error
	DeleteAdvert(id uint) error
	NewAdvertsList() *AdvertsList
	FillAdverts(ads *AdvertsList, imageService *services.ImageService)
}

type AdvertUseCase interface {
	CreateAdvert(advert *Advert) error
	GetAllAdverts() ([]*Advert, error)
	GetAdvertById(id uint) (*Advert, error)
	UpdateAdvert(advert *Advert) error
	DeleteAdvert(id uint) error
}
