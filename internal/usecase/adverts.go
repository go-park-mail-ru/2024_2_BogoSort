package usecase

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/google/uuid"
)

type AdvertUseCase interface {
	// GetAdverts возвращает массив объявлений в соответствии с offset и limit
	GetAdverts(limit, offset int) ([]*dto.AdvertResponse, error)

	// GetAdvertsByUserId возвращает массив объявлений в соответствии с userId
	GetAdvertsByUserId(userId uuid.UUID) ([]*dto.AdvertResponse, error)

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
	AddAdvert(advert *dto.AdvertRequest, userId uuid.UUID) (*dto.AdvertResponse, error)

	// UpdateAdvert обновляет объявление
	// Возможные ошибки:
	// ErrAdvertBadRequest - некорректные данные для обновления объявления
	// ErrAdvertNotFound - объявление не найдено
	// ErrForbidden - нет прав на обновление объявления
	UpdateAdvert(advert *dto.AdvertRequest, userId uuid.UUID, advertId uuid.UUID) error

	// UpdateAdvertStatus обновляет статус объявления
	// Возможные ошибки:
	// ErrAdvertBadRequest - некорректные данные для обновления статуса объявления
	// ErrAdvertNotFound - объявление не найдено
	// ErrForbidden - нет прав на обновление статуса объявления
	UpdateAdvertStatus(advertId uuid.UUID, status string, userId uuid.UUID) error

	// DeleteAdvertById удаляет объявление по Id
	// Возможные ошибки:
	// ErrAdvertNotFound - объявление не найдено
	// ErrForbidden - нет прав на удаление объявления
	DeleteAdvertById(advertId uuid.UUID, userId uuid.UUID) error

	// GetAdvertsByCategoryId возвращает массив объявлений по categoryId
	GetAdvertsByCategoryId(categoryId uuid.UUID) ([]*dto.AdvertResponse, error)

	// UploadImage загружает изображение в объявление
	// Возможные ошибки:
	// ErrAdvertNotFound - объявление не найдено
	// ErrForbidden - нет прав на загрузку изображения
	UploadImage(advertId uuid.UUID, imageId uuid.UUID, userId uuid.UUID) error
}
