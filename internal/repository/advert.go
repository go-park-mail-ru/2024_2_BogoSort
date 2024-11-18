package repository

import (
	"errors"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/google/uuid"
)

type AdvertRepository interface {
	// Get возвращает массив объявлений в соответствии с offset и limit
	Get(limit, offset int, userId uuid.UUID) ([]*entity.Advert, error)

	// GetBySellerId возвращает массив объявлений в соответствии с sellerId
	GetBySellerId(sellerId, userId uuid.UUID) ([]*entity.Advert, error)

	// GetByCartId возвращает массив объявлений, которые находятся в корзине
	GetByCartId(cartId uuid.UUID, userId uuid.UUID) ([]*entity.Advert, error)

	// GetByCategoryId возвращает массив объявлений по categoryId
	GetByCategoryId(categoryId, userId uuid.UUID) ([]*entity.Advert, error)

	// GetById возвращает объявление по его идентификатору
	// Если объявление не найдено, возвращает ErrAdvertNotFound
	GetById(advertId, userId uuid.UUID) (*entity.Advert, error)

	// GetSavedByUserId возвращает массив объявлений, которые находятся в сохраненных
	GetSavedByUserId(userId uuid.UUID) ([]*entity.Advert, error)

	// Add добавляет объявление
	// Возможные ошибки:
	// ErrAdvertBadRequest - некорректные данные для создания объявления
	// ErrAdvertAlreadyExists - объявление уже существует
	Add(advert *entity.Advert) (*entity.Advert, error)

	// AddToSaved добавляет объявление в сохраненные
	AddToSaved(userId, advertId uuid.UUID) error

	// DeleteFromSaved удаляет объявление из сохраненных
	DeleteFromSaved(userId, advertId uuid.UUID) error

	// Update обновляет объявление
	// Возможные ошибки:
	// ErrAdvertBadRequest - некорректные данные для создания объявления
	// ErrAdvertNotFound - объявление не найдено
	Update(advert *entity.Advert) error

	// DeleteById удаляет объявление по Id
	// Возможные ошибки:
	// ErrAdvertNotFound - объявление не найдено
	DeleteById(advertId uuid.UUID) error

	// UpdateStatus обновляет статус объявления
	// Возможные ошибки:
	// ErrAdvertBadRequest - некорректные данные для создания объявления
	// ErrAdvertNotFound - объявление не найдено
	UpdateStatus(tx pgx.Tx, advertId uuid.UUID, status entity.AdvertStatus) error

	// UploadImage загружает изображение в объявление
	UploadImage(advertId uuid.UUID, imageId uuid.UUID) error

	// AddViewed добавляет просмотренное объявление
	AddViewed(userId, advertId uuid.UUID) error

	// BeginTransaction начинает транзакцию
	BeginTransaction() (pgx.Tx, error)
}

var (
	ErrAdvertNotFound      = errors.New("объявление не найдено")
	ErrAdvertBadRequest    = errors.New("некорректные данные для создания объявления")
	ErrAdvertAlreadyExists = errors.New("объявление уже существует")
)
