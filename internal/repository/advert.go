package repository

import (
	"errors"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/google/uuid"
)

type AdvertRepository interface {
	// GetAdverts возвращает массив объявлений в соответствии с offset и limit
	GetAdverts(limit, offset int) ([]*entity.Advert, error)

	// GetAdvertsBySellerId возвращает массив объявлений в соответствии с sellerId
	GetAdvertsBySellerId(sellerId uuid.UUID) ([]*entity.Advert, error)

	// GetSavedAdvertsByUserId возвращает массив сохраненных объявлений в соответствии userId
	GetSavedAdvertsByUserId(userId uuid.UUID) ([]*entity.Advert, error)

	// GetAdvertsByCartId возвращает массив объявлений, которые находятся в корзине
	GetAdvertsByCartId(cartId uuid.UUID) ([]*entity.Advert, error)

	// GetAdvertsByCategoryId возвращает массив объявлений по categoryId
	GetAdvertsByCategoryId(categoryId uuid.UUID) ([]*entity.Advert, error)

	// GetAdvertById возвращает объявление по его идентификатору
	// Если объявление не найдено, возвращает ErrAdvertNotFound
	GetAdvertById(advertId uuid.UUID) (*entity.Advert, error)

	// AddAdvert добавляет объявление
	// Возможные ошибки:
	// ErrAdvertBadRequest - некорректные данные для создания объявления
	// ErrAdvertAlreadyExists - объявление уже существует
	AddAdvert(advert *entity.Advert) (*entity.Advert, error)

	// UpdateAdvert обновляет объявление
	// Возможные ошибки:
	// ErrAdvertBadRequest - некорректные данные для создания объявления
	// ErrAdvertNotFound - объявление не найдено
	UpdateAdvert(advert *entity.Advert) error

	// DeleteAdvertById удаляет объявление по Id
	// Возможные ошибки:
	// ErrAdvertNotFound - объявление не найдено
	DeleteAdvertById(advertId uuid.UUID) error

	// UpdateAdvertStatus обновляет статус объявления
	// Возможные ошибки:
	// ErrAdvertBadRequest - некорректные данные для создания объявления
	// ErrAdvertNotFound - объявление не найдено
	UpdateAdvertStatus(advertId uuid.UUID, status string) error

	// UploadImage загружает изображение в объявление
	UploadImage(advertId uuid.UUID, imageId uuid.UUID) error
}

var (
	ErrAdvertNotFound      = errors.New("объявление не найдено")
	ErrAdvertBadRequest    = errors.New("некорректные данные для создания объявления")
	ErrAdvertAlreadyExists = errors.New("объявление уже существует")
)
