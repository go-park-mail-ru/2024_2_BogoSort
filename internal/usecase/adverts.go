package usecase

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/google/uuid"
)	

type AdvertUseCase interface {
	// Get возвращает массив объявлений в соответствии с offset и limit
	Get(limit, offset int, userId uuid.UUID) ([]*dto.PreviewAdvertCard, error)

	// GetByUserId возвращает массив объявлений в соответствии с userId
	GetByUserId(userId uuid.UUID) ([]*dto.MyPreviewAdvertCard, error)

	// GetByCartId возвращает массив объявлений, которые находятся в корзине
	GetByCartId(cartId, userId uuid.UUID) ([]*dto.PreviewAdvertCard, error)

	// GetById возвращает объявление по его идентификатору
	// Если объявление не найдено, возвращает ErrAdvertNotFound
	GetById(advertId, userId uuid.UUID) (*dto.AdvertCard, error)

	// GetSavedByUserId возвращает массив объявлений, которые находятся в сохраненных
	GetSavedByUserId(userId uuid.UUID) ([]*dto.PreviewAdvertCard, error)

	// Add добавляет объявление
	// Возможные ошибки:
	// ErrAdvertBadRequest - некорректные данные для создания объявления
	// ErrAdvertAlreadyExists - объявление уже существует
	Add(advert *dto.AdvertRequest, userId uuid.UUID) (*dto.Advert, error)

	// Update обновляет объявление
	// Возможные ошибки:
	// ErrAdvertBadRequest - некорректные данные для обновления объявления
	// ErrAdvertNotFound - объявление не найдено
	// ErrForbidden - нет прав на обновление объявления
	Update(advert *dto.AdvertRequest, userId, advertId uuid.UUID) error

	// UpdateStatus обновляет статус объявления
	// Возможные ошибки:
	// ErrAdvertBadRequest - некорректные данные для обновления статуса объявления
	// ErrAdvertNotFound - объявление не найдено
	// ErrForbidden - нет прав на обновление статуса объявления
	UpdateStatus(advertId, userId uuid.UUID, status dto.AdvertStatus) error

	// DeleteById удаляет объявление по Id
	// Возможные ошибки:
	// ErrAdvertNotFound - объявление не найдено
	// ErrForbidden - нет прав на удаление объявления
	DeleteById(advertId uuid.UUID, userId uuid.UUID) error

	// GetByCategoryId возвращает массив объявлений по categoryId
	GetByCategoryId(categoryId, userId uuid.UUID) ([]*dto.PreviewAdvertCard, error)

	// UploadImage загружает изображение в объявление
	// Возможные ошибки:
	// ErrAdvertNotFound - объявление не найдено
	// ErrForbidden - нет прав на загрузку изображения
	UploadImage(advertId uuid.UUID, imageId uuid.UUID, userId uuid.UUID) error

	// AddToSaved добавляет объявление в сохраненные
	AddToSaved(advertId, userId uuid.UUID) error

	// RemoveFromSaved удаляет объявление из сохраненных
	RemoveFromSaved(advertId, userId uuid.UUID) error
}
