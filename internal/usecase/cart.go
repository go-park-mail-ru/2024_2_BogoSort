package usecase

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/google/uuid"
)

type Cart interface {
	// AddAdvert добавляет товар в корзину юзера по его ID
	AddAdvert(userID uuid.UUID, AdvertID uuid.UUID) error
	// GetByID возвращает корзину по ID корзины
	GetById(cartID uuid.UUID) (dto.Cart, error)
	// GetByUserID возвращает корзину по ID юзера
	GetByUserId(userID uuid.UUID) (dto.Cart, error)
	// DeleteAdvert удаляет товар из корзины по ID корзины и ID товара
	DeleteAdvert(cartID uuid.UUID, AdvertID uuid.UUID) error
	// CheckExists проверяет, существует ли корзина для пользователя
	CheckExists(userID uuid.UUID) (uuid.UUID, error)
}
