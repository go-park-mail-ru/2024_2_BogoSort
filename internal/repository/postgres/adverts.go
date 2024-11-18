package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type AdvertDB struct {
	DB      DBExecutor
	logger  *zap.Logger
	ctx     context.Context
	timeout time.Duration
}

const (
	insertAdvertQuery = `
		INSERT INTO advert (title, description, price, location, has_delivery, category_id, seller_id, image_id, status) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		RETURNING id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status`

	selectAdvertsQuery = `
		SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status, created_at, updated_at
		FROM advert
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	selectSavedAdvertsByUserIdQuery = `
		SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status, created_at, updated_at
		FROM advert
		WHERE id IN (SELECT advert_id FROM saved_advert WHERE user_id = $1)
		ORDER BY created_at DESC`

	selectAdvertsByUserIdQuery = `
		SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status, created_at, updated_at
		FROM advert
		WHERE seller_id = $1
		ORDER BY created_at DESC`

	selectAdvertsByCartIdQuery = `
		SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status, created_at, updated_at
		FROM advert
		WHERE id IN (SELECT advert_id FROM cart_advert WHERE cart_id = $1)
		ORDER BY created_at DESC`

	selectAdvertByIdQuery = `
		SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status, created_at, updated_at
		FROM advert
		WHERE id = $1
		ORDER BY created_at DESC`

	updateAdvertQuery = `
		UPDATE advert
		SET title = $1, description = $2, price = $3, location = $4, has_delivery = $5,
				category_id = $6, status = $7
		WHERE id = $8`

	deleteAdvertByIdQuery = `DELETE FROM advert WHERE id = $1`

	updateAdvertStatusQuery = `
		UPDATE advert
		SET status = $1
		WHERE id = $2`

	selectAdvertsByCategoryIdQuery = `
		SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status, created_at, updated_at
		FROM advert
		WHERE category_id = $1
		ORDER BY created_at DESC`

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
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewAdvertRepository(db *pgxpool.Pool, logger *zap.Logger, ctx context.Context, timeout time.Duration) (repository.AdvertRepository, error) {
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	return &AdvertDB{
		DB:      db,
		logger:  logger,
		ctx:     ctx,
		timeout: timeout,
	}, nil
}

func (r *AdvertDB) Add(a *entity.Advert) (*entity.Advert, error) {
	var dbAdvert AdvertRepoModel

	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	err := r.DB.QueryRow(ctx, insertAdvertQuery,
		a.Title,
		a.Description,
		a.Price,
		a.Location,
		a.HasDelivery,
		a.CategoryId,
		a.SellerId,
		a.ImageURL,
		string(a.Status)).Scan(
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
		return nil, entity.PSQLWrap(err)
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

func (r *AdvertDB) Get(limit, offset int) ([]*entity.Advert, error) {
	var adverts []*entity.Advert

	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	rows, err := r.DB.Query(ctx, selectAdvertsQuery, limit, offset)
	if err != nil {
		r.logger.Error("failed to execute query", zap.Error(err))
		return nil, entity.PSQLWrap(err)
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
			&dbAdvert.CreatedAt,
			&dbAdvert.UpdatedAt,
		); err != nil {
			r.logger.Error("failed to scan row", zap.Error(err))
			return nil, entity.PSQLWrap(err)
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
			CreatedAt:   dbAdvert.CreatedAt,
			UpdatedAt:   dbAdvert.UpdatedAt,
		})
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("error iterating over rows", zap.Error(err))
		return nil, entity.PSQLWrap(err)
	}

	return adverts, nil
}

func (r *AdvertDB) GetByCategoryId(categoryId uuid.UUID) ([]*entity.Advert, error) {
	var adverts []*entity.Advert

	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	rows, err := r.DB.Query(ctx, selectAdvertsByCategoryIdQuery, categoryId)
	if err != nil {
		r.logger.Error("failed to execute query", zap.Error(err), zap.String("category_id", categoryId.String()))
		return nil, entity.PSQLWrap(err)
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
			&dbAdvert.CreatedAt,
			&dbAdvert.UpdatedAt,
		); err != nil {
			r.logger.Error("failed to scan row", zap.Error(err), zap.String("category_id", categoryId.String()))
			return nil, entity.PSQLWrap(err)
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
			CreatedAt:   dbAdvert.CreatedAt,
			UpdatedAt:   dbAdvert.UpdatedAt,
		})
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("error iterating over rows", zap.Error(err), zap.String("category_id", categoryId.String()))
		return nil, entity.PSQLWrap(err)
	}

	return adverts, nil
}

func (r *AdvertDB) GetBySellerId(sellerId uuid.UUID) ([]*entity.Advert, error) {
	var adverts []*entity.Advert

	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	rows, err := r.DB.Query(ctx, selectAdvertsByUserIdQuery, sellerId)
	if err != nil {
		r.logger.Error("failed to execute query", zap.Error(err), zap.String("seller_id", sellerId.String()))
		return nil, entity.PSQLWrap(err)
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
			&dbAdvert.CreatedAt,
			&dbAdvert.UpdatedAt,
		); err != nil {
			r.logger.Error("failed to scan row", zap.Error(err), zap.String("seller_id", sellerId.String()))
			return nil, entity.PSQLWrap(err)
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
			CreatedAt:   dbAdvert.CreatedAt,
			UpdatedAt:   dbAdvert.UpdatedAt,
		})
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("error iterating over rows", zap.Error(err), zap.String("seller_id", sellerId.String()))
		return nil, entity.PSQLWrap(err)
	}

	return adverts, nil
}

func (r *AdvertDB) GetByCartId(cartId uuid.UUID) ([]*entity.Advert, error) {
	var adverts []*entity.Advert

	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	rows, err := r.DB.Query(ctx, selectAdvertsByCartIdQuery, cartId)
	if err != nil {
		r.logger.Error("failed to execute query", zap.Error(err), zap.String("user_id", cartId.String()))
		return nil, entity.PSQLWrap(err)
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
			&dbAdvert.CreatedAt,
			&dbAdvert.UpdatedAt,
		); err != nil {
			r.logger.Error("failed to scan row", zap.Error(err), zap.String("cart_id", cartId.String()))
			return nil, entity.PSQLWrap(err)
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
			CreatedAt:   dbAdvert.CreatedAt,
			UpdatedAt:   dbAdvert.UpdatedAt,
		})
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("error iterating over rows", zap.Error(err), zap.String("user_id", cartId.String()))
		return nil, entity.PSQLWrap(err)
	}

	return adverts, nil
}

func (r *AdvertDB) BeginTransaction() (pgx.Tx, error) {
	tx, err := r.DB.Begin(r.ctx)
	if err != nil {
		r.logger.Error("failed to begin transaction", zap.Error(err))
		return nil, err
	}
	return tx, nil
}

func (r *AdvertDB) GetById(advertId uuid.UUID) (*entity.Advert, error) {
	var dbAdvert AdvertRepoModel

	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

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
		&dbAdvert.CreatedAt,
		&dbAdvert.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Error("advert not found", zap.Error(err), zap.String("advert_id", advertId.String()))
			return nil, entity.PSQLWrap(repository.ErrAdvertNotFound)
		}
		r.logger.Error("failed to query advert by id", zap.Error(err), zap.String("advert_id", advertId.String()))
		return nil, entity.PSQLWrap(err)
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
		CreatedAt:   dbAdvert.CreatedAt,
		UpdatedAt:   dbAdvert.UpdatedAt,
	}, nil
}

func (r *AdvertDB) GetSavedByUserId(userId uuid.UUID) ([]*entity.Advert, error) {
	var adverts []*entity.Advert

	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	rows, err := r.DB.Query(ctx, selectSavedAdvertsByUserIdQuery, userId)
	if err != nil {
		r.logger.Error("failed to execute query", zap.Error(err), zap.String("user_id", userId.String()))
		return nil, entity.PSQLWrap(err)
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
			&dbAdvert.CreatedAt,
			&dbAdvert.UpdatedAt,
		); err != nil {
			r.logger.Error("failed to scan row", zap.Error(err), zap.String("user_id", userId.String()))
			return nil, entity.PSQLWrap(err)
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
			CreatedAt:   dbAdvert.CreatedAt,
			UpdatedAt:   dbAdvert.UpdatedAt,
		})
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("error iterating over rows", zap.Error(err), zap.String("user_id", userId.String()))
		return nil, entity.PSQLWrap(err)
	}

	return adverts, nil
}

func (r *AdvertDB) Update(advert *entity.Advert) error {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Minute)
	defer cancel()

	result, err := r.DB.Exec(ctx, updateAdvertQuery,
		advert.Title,
		advert.Description,
		advert.Price,
		advert.Location,
		advert.HasDelivery,
		advert.CategoryId,
		advert.Status,
		advert.ID,
	)
	if err != nil {
		r.logger.Error("failed to update advert", zap.Error(err), zap.String("advert_id", advert.ID.String()))
		return entity.PSQLWrap(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		r.logger.Error("advert not found", zap.String("advert_id", advert.ID.String()))
		return entity.PSQLWrap(repository.ErrAdvertNotFound)
	}

	return nil
}

func (r *AdvertDB) DeleteById(advertId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	result, err := r.DB.Exec(ctx, deleteAdvertByIdQuery, advertId)
	if err != nil {
		r.logger.Error("failed to delete advert", zap.Error(err), zap.String("advert_id", advertId.String()))
		return entity.PSQLWrap(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		r.logger.Error("advert not found", zap.String("advert_id", advertId.String()))
		return entity.PSQLWrap(repository.ErrAdvertNotFound)
	}

	return nil
}

func (r *AdvertDB) UpdateStatus(tx pgx.Tx, advertId uuid.UUID, status entity.AdvertStatus) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	result, err := tx.Exec(ctx, updateAdvertStatusQuery, status, advertId)
	if err != nil {
		r.logger.Error("failed to update advert status", zap.Error(err), zap.String("advert_id", advertId.String()))
		return entity.PSQLWrap(err)
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
		return entity.PSQLWrap(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		r.logger.Error("advert not found", zap.String("advert_id", advertId.String()))
		return entity.PSQLWrap(repository.ErrAdvertNotFound)
	}

	return nil
}
