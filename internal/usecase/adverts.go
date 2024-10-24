package usecase

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/google/uuid"
)

type AdvertUseCase interface {
	// GetAdverts возвращает массив объявлений в соответствии с offset и limit
	GetAdverts(limit, offset int) ([]*dto.Advert, error)

	// GetAdvertsByUserId возвращает массив объявлений в соответствии с userId
	GetAdvertsByUserId(userId uuid.UUID) ([]*dto.Advert, error)

	// GetSavedAdvertsByUserId возвращает массив сохраненных объявлений в соответствии userId
	GetSavedAdvertsByUserId(userId uuid.UUID) ([]*dto.Advert, error)

	// GetAdvertsByCartId возвращает массив объявлений, которые находятся в корзине
	GetAdvertsByCartId(cartId uuid.UUID) ([]*dto.Advert, error)

	// GetAdvertById возвращает объявление по его идентификатору
	// Если объявление не найдено, возвращает ErrAdvertNotFound
	GetAdvertById(advertId uuid.UUID) (*dto.Advert, error)

	// AddAdvert добавляет объявление
	// Возможные ошибки:
	// ErrAdvertBadRequest - некорректные данные для создания объявления
	// ErrAdvertAlreadyExists - объявление уже существует
	AddAdvert(advert *dto.Advert) (*dto.Advert, error)

	// UpdateAdvert обновляет объявление
	// Возможные ошибки:
	// ErrAdvertBadRequest - некорректные данные для создания объявления
	// ErrAdvertNotFound - объявление не найдено
	UpdateAdvert(advert *dto.Advert) error

	// DeleteAdvertById удаляет объявление по Id
	// Возможные ошибки:
	// ErrAdvertNotFound - объявление не найдено
	DeleteAdvertById(advertId uuid.UUID) error
}