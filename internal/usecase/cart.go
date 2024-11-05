package usecase

import (
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/google/uuid"
)

type Cart interface {
	// AddAdvertToUserCart добавляет товар в корзину юзера по его ID
	AddAdvertToUserCart(userID uuid.UUID, AdvertID uuid.UUID) error
	// GetAdvertsByCartID возвращает корзину по ID корзины
	GetCartByID(cartID uuid.UUID) (dto.Cart, error)
	// GetCartByUserID возвращает корзину по ID юзера
	GetCartByUserID(userID uuid.UUID) (dto.Cart, error)
	// DeleteAdvertFromCart удаляет товар из корзины по ID корзины и ID товара
	DeleteAdvertFromCart(cartID uuid.UUID, AdvertID uuid.UUID) error
	// CheckCartExists проверяет, существует ли корзина для пользователя
	CheckCartExists(userID uuid.UUID) (bool, error)
}
