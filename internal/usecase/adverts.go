package usecase

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/google/uuid"
)

type AdvertUseCase interface {
	// GetAdverts возвращает массив объявлений в соответствии с offset и limit
	GetAdverts(limit, offset int) ([]*dto.AdvertResponse, error)

	// GetAdvertsBySellerId возвращает массив объявлений в соответствии с sellerId
	GetAdvertsBySellerId(sellerId uuid.UUID) ([]*dto.AdvertResponse, error)

	// GetSavedAdvertsByUserId возвращает массив сохраненных объявлений в соответствии userId
	GetSavedAdvertsByUserId(userId uuid.UUID) ([]*dto.AdvertResponse, error)

	// GetAdvertsByCartId возвращает массив объявлений, которые находятся в корзине
	GetAdvertsByCartId(cartId uuid.UUID) ([]*dto.AdvertResponse, error)

	// GetAdvertById возвращает объявление по его идентификатору
	// Если объявление не найдено, возвращает ErrAdvertNotFound
	GetAdvertById(advertId uuid.UUID) (*dto.AdvertResponse, error)

	// AddAdvert добавляет объявление
	// Возможные ошибки:
	// ErrAdvertBadRequest - некорректные данные для создания объявления
	// ErrAdvertAlreadyExists - объявление уже существует
	AddAdvert(advert *dto.AdvertRequest) (*dto.AdvertResponse, error)

	// UpdateAdvert обновляет объявление
	// Возможные ошибки:
	// ErrAdvertBadRequest - некорректные данные для обновления объявления
	// ErrAdvertNotFound - объявление не найдено
	UpdateAdvert(advert *dto.AdvertRequest) error

	// UpdateAdvertStatus обновляет статус объявления
	// Возможные ошибки:
	// ErrAdvertBadRequest - некорректные данные для обновления статуса объявления
	// ErrAdvertNotFound - объявление не найдено
	UpdateAdvertStatus(advertId uuid.UUID, status string) error

	// DeleteAdvertById удаляет объявление по Id
	// Возможные ошибки:
	// ErrAdvertNotFound - объявление не найдено
	DeleteAdvertById(advertId uuid.UUID) error

	// GetAdvertsByCategoryId возвращает массив объявлений по categoryId
	GetAdvertsByCategoryId(categoryId uuid.UUID) ([]*dto.AdvertResponse, error)

	// UploadImage загружает изображение в объявление
	UploadImage(advertId uuid.UUID, imageId uuid.UUID) error
}
