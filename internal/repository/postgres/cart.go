package postgres

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/middleware"
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
	DB  DBExecutor
	ctx context.Context
}

func NewCartRepository(db *pgxpool.Pool, ctx context.Context) (repository.Cart, error) {
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	return &CartDB{
		DB:  db,
		ctx: ctx,
	}, nil
}

func (c *CartDB) GetByUserId(userID uuid.UUID) (entity.Cart, error) {
	var cart entity.Cart
	logger := middleware.GetLogger(c.ctx)
	logger.Info("getting cart by user id from db", zap.String("user_id", userID.String()))

	err := c.DB.QueryRow(c.ctx, queryGetCart, userID).Scan(&cart.ID, &cart.UserID, &cart.Status)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		logger.Error("cart not found", zap.String("user_id", userID.String()))
		return entity.Cart{}, repository.ErrCartNotFound
	case err != nil:
		logger.Error("error getting cart by user id", zap.String("user_id", userID.String()), zap.Error(err))
		return entity.Cart{}, entity.PSQLWrap(errors.New("error getting cart by user id"), err)
	}

	return cart, nil
}

func (c *CartDB) Create(userID uuid.UUID) (uuid.UUID, error) {
	var cartID uuid.UUID
	logger := middleware.GetLogger(c.ctx)
	logger.Info("creating cart in db", zap.String("user_id", userID.String()))

	err := c.DB.QueryRow(c.ctx, queryCreateCart, userID).Scan(&cartID)

	switch {
	case err != nil:
		logger.Error("error creating cart", zap.String("user_id", userID.String()), zap.Error(err))
		return uuid.Nil, entity.PSQLWrap(errors.New("error creating cart"), err)
	}
	return cartID, nil
}

func (c *CartDB) AddAdvert(cartID uuid.UUID, AdvertID uuid.UUID) error {
	logger := middleware.GetLogger(c.ctx)
	logger.Info("adding advert to cart in db", zap.String("cart_id", cartID.String()), zap.String("advert_id", AdvertID.String()))

	_, err := c.DB.Exec(c.ctx, queryAddAdvertToCart, cartID, AdvertID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.ErrCartNotFound
	case err != nil:
		logger.Error("error adding Advert to cart", zap.String("cart_id", cartID.String()), zap.String("Advert_id", AdvertID.String()), zap.Error(err))
		return entity.PSQLWrap(errors.New("error adding Advert to cart"), err)
	}
	return nil
}

func (c *CartDB) DeleteAdvert(cartID uuid.UUID, AdvertID uuid.UUID) error {
	logger := middleware.GetLogger(c.ctx)
	logger.Info("deleting advert from cart in db", zap.String("cart_id", cartID.String()), zap.String("Advert_id", AdvertID.String()))

	cmdTag, err := c.DB.Exec(c.ctx, queryDeleteAdvertFromCart, cartID, AdvertID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.ErrCartOrAdvertNotFound
	case err != nil:
		logger.Error("error deleting Advert from cart", zap.String("cart_id", cartID.String()), zap.String("Advert_id", AdvertID.String()), zap.Error(err))
		return entity.PSQLWrap(errors.New("error deleting Advert from cart"), err)
	}

	if cmdTag.RowsAffected() == 0 {
		return repository.ErrCartOrAdvertNotFound
	}
	return nil
}

func (c *CartDB) GetAdvertsByCartId(cartID uuid.UUID) ([]entity.Advert, error) {
	logger := middleware.GetLogger(c.ctx)
	logger.Info("getting adverts by cart id from db", zap.String("cart_id", cartID.String()))

	rows, err := c.DB.Query(c.ctx, queryGetAdvertsByCartID, cartID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, repository.ErrCartNotFound
	case err != nil:
		logger.Error("error getting adverts by cart id", zap.String("cart_id", cartID.String()), zap.Error(err))
		return nil, entity.PSQLWrap(errors.New("error getting adverts by cart id"), err)
	}
	defer rows.Close()

	var Adverts []entity.Advert
	for rows.Next() {
		var Advert entity.Advert
		if err := rows.Scan(&Advert.ID, &Advert.Title, &Advert.Description, &Advert.Price, &Advert.Location, &Advert.HasDelivery, &Advert.Status); err != nil {
			logger.Error("error scanning adverts", zap.String("cart_id", cartID.String()), zap.Error(err))
			return nil, entity.PSQLWrap(errors.New("error scanning adverts"), err)
		}
		Adverts = append(Adverts, Advert)
	}

	return Adverts, nil
}

func (c *CartDB) UpdateStatus(tx pgx.Tx, cartID uuid.UUID, status entity.CartStatus) error {
	logger := middleware.GetLogger(c.ctx)
	logger.Info("updating cart status in db", zap.String("cart_id", cartID.String()), zap.String("status", string(status)))

	_, err := tx.Exec(c.ctx, queryUpdateCartStatus, cartID, status)
	switch {
	case err != nil:
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Error("cart not found", zap.String("cart_id", cartID.String()))
			return repository.ErrCartNotFound
		}
		logger.Error("error updating cart status", zap.String("cart_id", cartID.String()), zap.Error(err))
		return entity.PSQLWrap(errors.New("error updating cart status"), err)
	}
	return nil
}

func (c *CartDB) GetById(cartID uuid.UUID) (entity.Cart, error) {
	var cart entity.Cart
	logger := middleware.GetLogger(c.ctx)
	logger.Info("getting cart by id from db", zap.String("cart_id", cartID.String()))

	err := c.DB.QueryRow(c.ctx, queryGetCartByID, cartID).Scan(&cart.ID, &cart.UserID, &cart.Status)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return entity.Cart{}, repository.ErrCartNotFound
	case err != nil:
		logger.Error("error getting cart by id", zap.String("cart_id", cartID.String()), zap.Error(err))
		return entity.Cart{}, entity.PSQLWrap(errors.New("error getting cart by id"), err)
	}
	return cart, nil
}
