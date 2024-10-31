package postgres

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"time"
)

type AdvertDB struct {
	DB      *pgxpool.Pool
	logger  *zap.Logger
	timeout time.Duration
	ctx     context.Context
}

const (
	insertAdvertQuery = `
		INSERT INTO advert (title, description, price, location, has_delivery, category_id, seller_id, image_id, status) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		RETURNING id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status`

	selectAdvertsQuery = `
		SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status
		FROM advert
		LIMIT $1 OFFSET $2`

	selectAdvertsByUserIdQuery = `
		SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status
		FROM advert
		WHERE seller_id = $1`

	selectSavedAdvertsByUserIdQuery = `
		SELECT a.id, a.title, a.description, a.price, a.location, a.has_delivery, a.category_id, a.seller_id, a.image_id, a.status
		FROM advert a
		INNER JOIN saved_advert sa ON sa.advert_id = a.id
		WHERE sa.user_id = $1`

	selectAdvertsByCartIdQuery = `
		SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status
		FROM advert
		WHERE id IN (SELECT advert_id FROM cart_advert WHERE cart_id = $1)`

	selectAdvertByIdQuery = `
		SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status
		FROM advert
		WHERE id = $1`

	updateAdvertQuery = `
		UPDATE advert
		SET title = $1, description = $2, price = $3, location = $4, has_delivery = $5,
				category_id = $6, seller_id = $7, image_id = $8, status = $9, updated_at = NOW()
		WHERE id = $10`

	deleteAdvertByIdQuery = `DELETE FROM advert WHERE id = $1`

	updateAdvertStatusQuery = `
		UPDATE advert
		SET status = $1, updated_at = NOW()
		WHERE id = $2`

	selectAdvertsByCategoryIdQuery = `
		SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status
		FROM advert
		WHERE category_id = $1`

	uploadImageQuery = `
		UPDATE advert
		SET image_id = $1
		WHERE id = $2`
)

type AdvertRepoModel struct {
    ID          uuid.UUID
    SellerId    uuid.UUID
    CategoryId  uuid.UUID
    Title       string
    Description string
    Price       uint
    ImageURL    uuid.NullUUID
    Status      string
    HasDelivery bool
    Location    string
}

func NewAdvertRepository(db *pgxpool.Pool, logger *zap.Logger, timeout time.Duration, ctx context.Context) (repository.AdvertRepository, error) {
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	return &AdvertDB{
		DB:      db,
		logger:  logger,
		timeout: timeout,
		ctx:     ctx,
	}, nil
}

func (r *AdvertDB) AddAdvert(a *entity.Advert) (*entity.Advert, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	var dbAdvert AdvertRepoModel

	err := r.DB.QueryRow(ctx, insertAdvertQuery,
		a.Title,
		a.Description,
		a.Price,
		a.Location,
		a.HasDelivery,
		a.CategoryId,
		a.SellerId,
		a.ImageURL,
		a.Status).Scan(
		&dbAdvert.ID,
		&dbAdvert.Title,
		&dbAdvert.Description,
		&dbAdvert.Price,
		&dbAdvert.Location,
		&dbAdvert.HasDelivery,
		&dbAdvert.CategoryId,
		&dbAdvert.SellerId,
		&dbAdvert.ImageURL,
		&dbAdvert.Status,
	)

	if err != nil {
		r.logger.Error("error adding advert", zap.Error(err))
		return nil, entity.PSQLQueryErr("AddAdvert", err)
	}

	r.logger.Info("advert added", zap.Any("advert", dbAdvert))

	return &entity.Advert{
		ID:          dbAdvert.ID,
		Title:       dbAdvert.Title,
		Description: dbAdvert.Description,
		Price:       dbAdvert.Price,
		Location:    dbAdvert.Location,
		HasDelivery: dbAdvert.HasDelivery,
		CategoryId:  dbAdvert.CategoryId,
		SellerId:    dbAdvert.SellerId,
		ImageURL:    dbAdvert.ImageURL,
		Status:      entity.AdvertStatus(dbAdvert.Status),
	}, nil
}

func (r *AdvertDB) GetAdverts(limit, offset int) ([]*entity.Advert, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	var adverts []*entity.Advert

	rows, err := r.DB.Query(ctx, selectAdvertsQuery, limit, offset)
	if err != nil {
		r.logger.Error("failed to execute query", zap.Error(err))
		return nil, entity.PSQLQueryErr("GetAdverts", err)
	}
	defer rows.Close()

	for rows.Next() {
		var dbAdvert AdvertRepoModel
		if err := rows.Scan(&dbAdvert.ID,
			&dbAdvert.Title,
			&dbAdvert.Description,
			&dbAdvert.Price,
			&dbAdvert.Location,
			&dbAdvert.HasDelivery,
			&dbAdvert.CategoryId,
			&dbAdvert.SellerId,
			&dbAdvert.ImageURL,
			&dbAdvert.Status,
		); err != nil {
			r.logger.Error("failed to scan row", zap.Error(err))
			return nil, entity.PSQLQueryErr("GetAdverts", err)
		}
		adverts = append(adverts, &entity.Advert{
			ID:          dbAdvert.ID,
			Title:       dbAdvert.Title,
			Description: dbAdvert.Description,
			Price:       dbAdvert.Price,
			Location:    dbAdvert.Location,
			HasDelivery: dbAdvert.HasDelivery,
			CategoryId:  dbAdvert.CategoryId,
			SellerId:    dbAdvert.SellerId,
			ImageURL:    dbAdvert.ImageURL,
			Status:      entity.AdvertStatus(dbAdvert.Status),
		})
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("error iterating over rows", zap.Error(err))
		return nil, entity.PSQLQueryErr("GetAdverts", err)
	}

	return adverts, nil
}

func (r *AdvertDB) GetAdvertsByCategoryId(categoryId uuid.UUID) ([]*entity.Advert, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	var adverts []*entity.Advert

	rows, err := r.DB.Query(ctx, selectAdvertsByCategoryIdQuery, categoryId)
	if err != nil {
		r.logger.Error("failed to execute query", zap.Error(err), zap.String("category_id", categoryId.String()))
		return nil, entity.PSQLQueryErr("GetAdvertsByCategoryId", err)
	}

	for rows.Next() {
		var dbAdvert AdvertRepoModel
		if err := rows.Scan(
			&dbAdvert.ID,
			&dbAdvert.Title,
			&dbAdvert.Description,
			&dbAdvert.Price,
			&dbAdvert.Location,
			&dbAdvert.HasDelivery,
			&dbAdvert.CategoryId,
			&dbAdvert.SellerId,
			&dbAdvert.ImageURL,
			&dbAdvert.Status,
		); err != nil {
			r.logger.Error("failed to scan row", zap.Error(err), zap.String("category_id", categoryId.String()))
			return nil, entity.PSQLQueryErr("GetAdvertsByCategoryId", err)
		}
		adverts = append(adverts, &entity.Advert{
			ID:          dbAdvert.ID,
			Title:       dbAdvert.Title,
			Description: dbAdvert.Description,
			Price:       dbAdvert.Price,
			Location:    dbAdvert.Location,
			HasDelivery: dbAdvert.HasDelivery,
			CategoryId:  dbAdvert.CategoryId,
			SellerId:    dbAdvert.SellerId,
			ImageURL:    dbAdvert.ImageURL,
			Status:      entity.AdvertStatus(dbAdvert.Status),
		})
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("error iterating over rows", zap.Error(err), zap.String("category_id", categoryId.String()))
		return nil, entity.PSQLQueryErr("GetAdvertsByCategoryId", err)
	}

	return adverts, nil
}

func (r *AdvertDB) GetAdvertsBySellerId(sellerId uuid.UUID) ([]*entity.Advert, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	var adverts []*entity.Advert

	rows, err := r.DB.Query(ctx, selectAdvertsByUserIdQuery, sellerId)
	if err != nil {
		r.logger.Error("failed to execute query", zap.Error(err), zap.String("seller_id", sellerId.String()))
		return nil, entity.PSQLQueryErr("GetAdvertsBySellerId", err)
	}
	defer rows.Close()

	for rows.Next() {
		var dbAdvert AdvertRepoModel
		if err := rows.Scan(
			&dbAdvert.ID,
			&dbAdvert.Title,
			&dbAdvert.Description,
			&dbAdvert.Price,
			&dbAdvert.Location,
			&dbAdvert.HasDelivery,
			&dbAdvert.CategoryId,
			&dbAdvert.SellerId,
			&dbAdvert.ImageURL,
			&dbAdvert.Status,
		); err != nil {
			r.logger.Error("failed to scan row", zap.Error(err), zap.String("seller_id", sellerId.String()))
			return nil, entity.PSQLQueryErr("GetAdvertsBySellerId", err)
		}
		adverts = append(adverts, &entity.Advert{
			ID:          dbAdvert.ID,
			Title:       dbAdvert.Title,
			Description: dbAdvert.Description,
			Price:       dbAdvert.Price,
			Location:    dbAdvert.Location,
			HasDelivery: dbAdvert.HasDelivery,
			CategoryId:  dbAdvert.CategoryId,
			SellerId:    dbAdvert.SellerId,
			ImageURL:    dbAdvert.ImageURL,
			Status:      entity.AdvertStatus(dbAdvert.Status),
		})
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("error iterating over rows", zap.Error(err), zap.String("seller_id", sellerId.String()))
		return nil, entity.PSQLQueryErr("GetAdvertsBySellerId", err)
	}

	return adverts, nil
}

func (r *AdvertDB) GetSavedAdvertsByUserId(userId uuid.UUID) ([]*entity.Advert, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	var adverts []*entity.Advert

	rows, err := r.DB.Query(ctx, selectSavedAdvertsByUserIdQuery, userId)
	if err != nil {
		r.logger.Error("failed to execute query", zap.Error(err), zap.String("user_id", userId.String()))
		return nil, entity.PSQLQueryErr("GetSavedAdvertsByUserId", err)
	}
	defer rows.Close()

	for rows.Next() {
		var dbAdvert AdvertRepoModel
		if err := rows.Scan(
			&dbAdvert.ID,
			&dbAdvert.Title,
			&dbAdvert.Description,
			&dbAdvert.Price,
			&dbAdvert.Location,
			&dbAdvert.HasDelivery,
			&dbAdvert.CategoryId,
			&dbAdvert.SellerId,
			&dbAdvert.ImageURL,
			&dbAdvert.Status,
		); err != nil {
			r.logger.Error("failed to scan row", zap.Error(err), zap.String("user_id", userId.String()))
			return nil, entity.PSQLQueryErr("GetSavedAdvertsByUserId", err)
		}
		adverts = append(adverts, &entity.Advert{
			ID:          dbAdvert.ID,
			Title:       dbAdvert.Title,
			Description: dbAdvert.Description,
			Price:       dbAdvert.Price,
			Location:    dbAdvert.Location,
			HasDelivery: dbAdvert.HasDelivery,
			CategoryId:  dbAdvert.CategoryId,
			SellerId:    dbAdvert.SellerId,
			ImageURL:    dbAdvert.ImageURL,
			Status:      entity.AdvertStatus(dbAdvert.Status),
		})
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("error iterating over rows", zap.Error(err), zap.String("user_id", userId.String()))
		return nil, entity.PSQLQueryErr("GetSavedAdvertsByUserId", err)
	}

	return adverts, nil
}

func (r *AdvertDB) GetAdvertsByCartId(cartId uuid.UUID) ([]*entity.Advert, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	var adverts []*entity.Advert

	rows, err := r.DB.Query(ctx, selectAdvertsByCartIdQuery, cartId)
	if err != nil {
		r.logger.Error("failed to execute query", zap.Error(err), zap.String("user_id", cartId.String()))
		return nil, entity.PSQLQueryErr("GetSavedAdvertsByUserId", err)
	}
	defer rows.Close()

	for rows.Next() {
		var dbAdvert AdvertRepoModel
		if err := rows.Scan(
			&dbAdvert.ID,
			&dbAdvert.Title,
			&dbAdvert.Description,
			&dbAdvert.Price,
			&dbAdvert.Location,
			&dbAdvert.HasDelivery,
			&dbAdvert.CategoryId,
			&dbAdvert.SellerId,
			&dbAdvert.ImageURL,
			&dbAdvert.Status,
		); err != nil {
			r.logger.Error("failed to scan row", zap.Error(err), zap.String("cart_id", cartId.String()))
			return nil, entity.PSQLQueryErr("GetSavedAdvertsByUserId", err)
		}
		adverts = append(adverts, &entity.Advert{
			ID:          dbAdvert.ID,
			Title:       dbAdvert.Title,
			Description: dbAdvert.Description,
			Price:       dbAdvert.Price,
			Location:    dbAdvert.Location,
			HasDelivery: dbAdvert.HasDelivery,
			CategoryId:  dbAdvert.CategoryId,
			SellerId:    dbAdvert.SellerId,
			ImageURL:    dbAdvert.ImageURL,
			Status:      entity.AdvertStatus(dbAdvert.Status),
		})
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("error iterating over rows", zap.Error(err), zap.String("user_id", cartId.String()))
		return nil, entity.PSQLQueryErr("GetSavedAdvertsByUserId", err)
	}

	return adverts, nil
}

func (r *AdvertDB) GetAdvertById(advertId uuid.UUID) (*entity.Advert, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	var dbAdvert AdvertRepoModel

	err := r.DB.QueryRow(ctx, selectAdvertByIdQuery, advertId).Scan(
		&dbAdvert.ID,
		&dbAdvert.Title,
		&dbAdvert.Description,
		&dbAdvert.Price,
		&dbAdvert.Location,
		&dbAdvert.HasDelivery,
		&dbAdvert.CategoryId,
		&dbAdvert.SellerId,
		&dbAdvert.ImageURL,
		&dbAdvert.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Error("advert not found", zap.Error(err), zap.String("advert_id", advertId.String()))
			return nil, entity.PSQLWrap(repository.ErrAdvertNotFound)
		}
		r.logger.Error("failed to query advert by id", zap.Error(err), zap.String("advert_id", advertId.String()))
		return nil, entity.PSQLQueryErr("GetAdvertById", err)
	}

	return &entity.Advert{
		ID:          dbAdvert.ID,
		Title:       dbAdvert.Title,
		Description: dbAdvert.Description,
		Price:       dbAdvert.Price,
		Location:    dbAdvert.Location,
		HasDelivery: dbAdvert.HasDelivery,
		CategoryId:  dbAdvert.CategoryId,
		SellerId:    dbAdvert.SellerId,
		ImageURL:    dbAdvert.ImageURL,
		Status:      entity.AdvertStatus(dbAdvert.Status),
	}, nil
}

func (r *AdvertDB) UpdateAdvert(advert *entity.Advert) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	result, err := r.DB.Exec(ctx, updateAdvertQuery,
		advert.Title,
		advert.Description,
		advert.Price,
			advert.Location,
			advert.HasDelivery,
			advert.CategoryId,
			advert.SellerId,
			advert.ImageURL,
			advert.Status,
			advert.ID,
	)
	if err != nil {
		r.logger.Error("failed to update advert", zap.Error(err), zap.String("advert_id", advert.ID.String()))
		return entity.PSQLQueryErr("UpdateAdvert", err)
	}

	rowsAffected:= result.RowsAffected()
	if rowsAffected == 0 {
		r.logger.Error("advert not found", zap.String("advert_id", advert.ID.String()))
		return entity.PSQLWrap(repository.ErrAdvertNotFound)
	}

	return nil
}

func (r *AdvertDB) DeleteAdvertById(advertId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	result, err := r.DB.Exec(ctx, deleteAdvertByIdQuery, advertId)
	if err != nil {
		r.logger.Error("failed to delete advert", zap.Error(err), zap.String("advert_id", advertId.String()))
		return entity.PSQLQueryErr("DeleteAdvertById", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		r.logger.Error("advert not found", zap.String("advert_id", advertId.String()))
		return entity.PSQLWrap(repository.ErrAdvertNotFound)
	}

	return nil
}

func (r *AdvertDB) UpdateAdvertStatus(advertId uuid.UUID, status string) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	result, err := r.DB.Exec(ctx, updateAdvertStatusQuery, status, advertId)
	if err != nil {
		r.logger.Error("failed to update advert status", zap.Error(err), zap.String("advert_id", advertId.String()))
		return entity.PSQLQueryErr("UpdateAdvertStatus", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		r.logger.Error("advert not found", zap.String("advert_id", advertId.String()))
		return entity.PSQLWrap(repository.ErrAdvertNotFound)
	}

	return nil
}

func (r *AdvertDB) UploadImage(advertId uuid.UUID, imageId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	result, err := r.DB.Exec(ctx, uploadImageQuery, imageId, advertId)
	if err != nil {
		r.logger.Error("failed to upload image", zap.Error(err), zap.String("advert_id", advertId.String()))
		return entity.PSQLQueryErr("UploadImage", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		r.logger.Error("advert not found", zap.String("advert_id", advertId.String()))
		return entity.PSQLWrap(repository.ErrAdvertNotFound)
	}

	return nil
}
