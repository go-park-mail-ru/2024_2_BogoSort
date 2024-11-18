package usecase

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/google/uuid"
)	

type AdvertUseCase interface {
	// Get возвращает массив объявлений в соответствии с offset и limit
	Get(limit, offset int) ([]*dto.AdvertResponse, error)

	// GetByUserId возвращает массив объявлений в соответствии с userId
	GetByUserId(userId uuid.UUID) ([]*dto.AdvertResponse, error)

	// GetByCartId возвращает массив объявлений, которые находятся в корзине
	GetByCartId(cartId uuid.UUID) ([]*dto.AdvertResponse, error)

	// GetById возвращает объявление по его идентификатору
	// Если объявление не найдено, возвращает ErrAdvertNotFound
	GetById(advertId uuid.UUID) (*dto.AdvertResponse, error)

	// GetSavedByUserId возвращает массив объявлений, которые находятся в сохраненных
	GetSavedByUserId(userId uuid.UUID) ([]*dto.AdvertResponse, error)

	// Add добавляет объявление
	// Возможные ошибки:
	// ErrAdvertBadRequest - некорректные данные для создания объявления
	// ErrAdvertAlreadyExists - объявление уже существует
	Add(advert *dto.AdvertRequest, userId uuid.UUID) (*dto.AdvertResponse, error)

	// Update обновляет объявление
	// Возможные ошибки:
	// ErrAdvertBadRequest - некорректные данные для обновления объявления
	// ErrAdvertNotFound - объявление не найдено
	// ErrForbidden - нет прав на обновление объявления
	Update(advert *dto.AdvertRequest, userId uuid.UUID, advertId uuid.UUID) error

	// UpdateStatus обновляет статус объявления
	// Возможные ошибки:
	// ErrAdvertBadRequest - некорректные данные для обновления статуса объявления
	// ErrAdvertNotFound - объявление не найдено
	// ErrForbidden - нет прав на обновление статуса объявления
	UpdateStatus(advertId uuid.UUID, status dto.AdvertStatus, userId uuid.UUID) error

	// DeleteById удаляет объявление по Id
	// Возможные ошибки:
	// ErrAdvertNotFound - объявление не найдено
	// ErrForbidden - нет прав на удаление объявления
	DeleteById(advertId uuid.UUID, userId uuid.UUID) error

	// GetByCategoryId возвращает массив объявлений по categoryId
	GetByCategoryId(categoryId uuid.UUID) ([]*dto.AdvertResponse, error)

	// UploadImage загружает изображение в объявление
	// Возможные ошибки:
	// ErrAdvertNotFound - объявление не найдено
	// ErrForbidden - нет прав на загрузку изображения
	UploadImage(advertId uuid.UUID, imageId uuid.UUID, userId uuid.UUID) error
}
