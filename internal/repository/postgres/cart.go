package postgres

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	queryCreateCart = `
		INSERT INTO cart (user_id) VALUES ($1) RETURNING id
	`
	queryAddAdvertToCart = `
		INSERT INTO cart_advert (cart_id, advert_id) VALUES ($1, $2)
	`
	queryDeleteAdvertFromCart = `
		DELETE FROM cart_advert WHERE cart_id = $1 AND advert_id = $2
	`
	queryGetAdvertsByCartID = `
		SELECT a.id, a.title, a.description, a.price, a.location, a.has_delivery, a.status
		FROM cart_advert ca
		JOIN advert a ON ca.advert_id = a.id
		WHERE ca.cart_id = $1
	`
	queryGetCart = `
		SELECT id, user_id, status FROM cart WHERE user_id = $1 AND status = 'active' LIMIT 1
	`
	queryUpdateCartStatus = `
		UPDATE cart SET status = $2 WHERE id = $1
	`
	queryGetCartByID = `
		SELECT id, user_id, status FROM cart WHERE id = $1
	`
)

type CartDB struct {
	DB     *pgxpool.Pool
	ctx    context.Context
	logger *zap.Logger
}

func NewCartRepository(db *pgxpool.Pool, ctx context.Context, logger *zap.Logger) (repository.Cart, error) {
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	return &CartDB{
		DB:     db,
		ctx:    ctx,
		logger: logger,
	}, nil
}

func (c *CartDB) GetCartByUserID(userID uuid.UUID) (entity.Cart, error) {
	var cart entity.Cart
	err := c.DB.QueryRow(c.ctx, queryGetCart, userID).Scan(&cart.ID, &cart.UserID, &cart.Status)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return entity.Cart{}, repository.ErrCartNotFound
	case err != nil:
		c.logger.Error("error getting cart by user id", zap.String("user_id", userID.String()), zap.Error(err))
		return entity.Cart{}, entity.PSQLWrap(errors.New("error getting cart by user id"), err)
	}

	return cart, nil
}

func (c *CartDB) CreateCart(userID uuid.UUID) (uuid.UUID, error) {
	var cartID uuid.UUID
	err := c.DB.QueryRow(c.ctx, queryCreateCart, userID).Scan(&cartID)

	switch {
	case err != nil:
		c.logger.Error("error creating cart", zap.String("user_id", userID.String()), zap.Error(err))
		return uuid.Nil, entity.PSQLWrap(errors.New("error creating cart"), err)
	}
	return cartID, nil
}

func (c *CartDB) AddAdvertToCart(cartID uuid.UUID, AdvertID uuid.UUID) error {
	_, err := c.DB.Exec(c.ctx, queryAddAdvertToCart, cartID, AdvertID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.ErrCartNotFound
	case err != nil:
		c.logger.Error("error adding Advert to cart", zap.String("cart_id", cartID.String()), zap.String("Advert_id", AdvertID.String()), zap.Error(err))
		return entity.PSQLWrap(errors.New("error adding Advert to cart"), err)
	}
	return nil
}

func (c *CartDB) DeleteAdvertFromCart(cartID uuid.UUID, AdvertID uuid.UUID) error {
	_, err := c.DB.Exec(c.ctx, queryDeleteAdvertFromCart, cartID, AdvertID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.ErrCartOrAdvertNotFound
	case err != nil:
		c.logger.Error("error deleting Advert from cart", zap.String("cart_id", cartID.String()), zap.String("Advert_id", AdvertID.String()), zap.Error(err))
		return entity.PSQLWrap(errors.New("error deleting Advert from cart"), err)
	}
	return nil
}

func (c *CartDB) GetAdvertsByCartID(cartID uuid.UUID) ([]entity.Advert, error) {
	rows, err := c.DB.Query(c.ctx, queryGetAdvertsByCartID, cartID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, repository.ErrCartNotFound
	case err != nil:
		c.logger.Error("error getting adverts by cart id", zap.String("cart_id", cartID.String()), zap.Error(err))
		return nil, entity.PSQLWrap(errors.New("error getting adverts by cart id"), err)
	}
	defer rows.Close()

	var Adverts []entity.Advert
	for rows.Next() {
		var Advert entity.Advert
		if err := rows.Scan(&Advert.ID, &Advert.Title, &Advert.Description, &Advert.Price, &Advert.Location, &Advert.HasDelivery, &Advert.Status); err != nil {
			c.logger.Error("error scanning adverts", zap.String("cart_id", cartID.String()), zap.Error(err))
			return nil, entity.PSQLWrap(errors.New("error scanning adverts"), err)
		}
		Adverts = append(Adverts, Advert)
	}

	return Adverts, nil
}

func (c *CartDB) UpdateCartStatus(tx pgx.Tx, cartID uuid.UUID, status entity.CartStatus) error {
	_, err := tx.Exec(c.ctx, queryUpdateCartStatus, cartID, status)
	switch {
	case err != nil:
		if errors.Is(err, pgx.ErrNoRows) {
			c.logger.Error("cart not found", zap.String("cart_id", cartID.String()))
			return repository.ErrCartNotFound
		}
		c.logger.Error("error updating cart status", zap.String("cart_id", cartID.String()), zap.Error(err))
		return entity.PSQLWrap(errors.New("error updating cart status"), err)
	}
	return nil
}

func (c *CartDB) GetCartByID(cartID uuid.UUID) (entity.Cart, error) {
	var cart entity.Cart
	err := c.DB.QueryRow(c.ctx, queryGetCartByID, cartID).Scan(&cart.ID, &cart.UserID, &cart.Status)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return entity.Cart{}, repository.ErrCartNotFound
	case err != nil:
		c.logger.Error("error getting cart by id", zap.String("cart_id", cartID.String()), zap.Error(err))
		return entity.Cart{}, entity.PSQLWrap(errors.New("error getting cart by id"), err)
	}
	return cart, nil
}
