package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/middleware"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type AdvertDB struct {
	DB      DBExecutor
	ctx     context.Context
	timeout time.Duration
}

const (
	insertAdvertQuery = `
		INSERT INTO advert (title, description, price, location, has_delivery, category_id, seller_id, status) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status`

	selectAdvertsQuery = `
		SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status, created_at, updated_at
		FROM advert
		WHERE status != 'inactive'
		ORDER BY 
    		CASE WHEN promoted_until > CURRENT_TIMESTAMP THEN 1 ELSE 0 END DESC, 
    		promoted_until DESC NULLS LAST, 
    		created_at DESC
		LIMIT $1 OFFSET $2;
`

	selectSavedAdvertsByUserIdQuery = `
		SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status, created_at, updated_at
		FROM advert
		WHERE id IN (SELECT advert_id FROM saved_advert WHERE user_id = $1)
		ORDER BY created_at DESC`

	selectAdvertsBySellerIdQuery = `
		SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status, created_at, updated_at
		FROM advert
		WHERE seller_id = $1 AND status != 'inactive'
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
		WHERE category_id = $1 AND status != 'inactive'
		ORDER BY created_at DESC`

	uploadImageQuery = `
		UPDATE advert
		SET image_id = $1
		WHERE id = $2`

	insertSavedAdvertQuery = `
		INSERT INTO saved_advert (user_id, advert_id)
		VALUES ($1, $2)
		RETURNING id, user_id, advert_id, created_at`

	deleteSavedAdvertQuery = `
		DELETE FROM saved_advert
		WHERE user_id = $1 AND advert_id = $2`

	selectSavedCountAndIsSavedQuery = `
		SELECT COUNT(*), EXISTS(SELECT 1 FROM saved_advert WHERE advert_id = $1 AND user_id = $2) 
		FROM saved_advert WHERE advert_id = $1`

	insertViewedAdvertQuery = `
		INSERT INTO viewed_advert (advert_id, user_id)
		VALUES ($1, $2)
		RETURNING id, user_id, advert_id, created_at`

	selectViewedCountAndIsViewedQuery = `
		SELECT COUNT(*), EXISTS(SELECT 1 FROM viewed_advert WHERE advert_id = $1 AND user_id = $2) 
		FROM viewed_advert WHERE advert_id = $1`

	checkIfExistsQuery = `
		SELECT EXISTS(SELECT 1 FROM advert WHERE id = $1)`

	searchAdvertsQuery = `
		SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status, created_at, updated_at
		FROM advert
		WHERE to_tsvector('russian', title || ' ' || description) @@ plainto_tsquery('russian', $1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	countAdvertsQuery = `SELECT COUNT(*) FROM advert`

	promoteAdvertQuery = `
	UPDATE advert
	SET promoted_until = CURRENT_TIMESTAMP + INTERVAL '7 days', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND status != 'inactive'
		RETURNING id, title, promoted_until;
	`
)

type AdvertRepoModel struct {
	ID          uuid.UUID
	SellerId    uuid.UUID
	CategoryId  uuid.UUID
	Title       string
	Description string
	Price       uint
	ImageId     uuid.UUID
	Status      string
	HasDelivery bool
	Location    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SavedAdvertRepoModel struct {
	ID        uuid.UUID
	AdvertId  uuid.UUID
	UserId    uuid.UUID
	CreatedAt time.Time
}

func NewAdvertRepository(db *pgxpool.Pool, ctx context.Context, timeout time.Duration) (repository.AdvertRepository, error) {
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	return &AdvertDB{
		DB:      db,
		ctx:     ctx,
		timeout: timeout,
	}, nil
}

func (r *AdvertDB) getSavedCount(advertId uuid.UUID, userId uuid.UUID) (int, bool, error) {
	var savedCount int
	isSaved := false

	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	logger := middleware.GetLogger(r.ctx)
	logger.Info("getting saved count from db", zap.String("advert_id", advertId.String()), zap.String("user_id", userId.String()))
	err := r.DB.QueryRow(ctx, selectSavedCountAndIsSavedQuery, advertId, userId).Scan(&savedCount, &isSaved)
	if err != nil {
		logger.Error("failed to execute query", zap.Error(err), zap.String("advert_id", advertId.String()), zap.String("user_id", userId.String()))
		return 0, false, err
	}

	return savedCount, isSaved, nil
}

func (r *AdvertDB) getViewedCount(advertId uuid.UUID, userId uuid.UUID) (int, bool, error) {
	var viewedCount int
	isViewed := false

	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	logger := middleware.GetLogger(r.ctx)
	logger.Info("getting viewed count from db", zap.String("advert_id", advertId.String()), zap.String("user_id", userId.String()))

	err := r.DB.QueryRow(ctx, selectViewedCountAndIsViewedQuery, advertId, userId).Scan(&viewedCount, &isViewed)
	if err != nil {
		logger.Error("failed to execute query", zap.Error(err), zap.String("advert_id", advertId.String()), zap.String("user_id", userId.String()))
		return 0, false, err
	}

	return viewedCount, isViewed, nil
}

func (r *AdvertDB) convertToEntityAdvert(dbAdvert AdvertRepoModel, userId uuid.UUID) *entity.Advert {
	savedCount, isSaved, err := r.getSavedCount(dbAdvert.ID, userId)
	logger := middleware.GetLogger(r.ctx)
	logger.Info("getting saved count", zap.String("advert_id", dbAdvert.ID.String()), zap.String("user_id", userId.String()))
	if err != nil {
		logger.Error("failed to get saved count", zap.Error(err), zap.String("advert_id", dbAdvert.ID.String()), zap.String("user_id", userId.String()))
		return nil
	}

	logger.Info("getting viewed count", zap.String("advert_id", dbAdvert.ID.String()), zap.String("user_id", userId.String()))
	viewedCount, isViewed, err := r.getViewedCount(dbAdvert.ID, userId)
	if err != nil {
		logger.Error("failed to get viewed count", zap.Error(err), zap.String("advert_id", dbAdvert.ID.String()), zap.String("user_id", userId.String()))
		return nil
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
		ImageId:     dbAdvert.ImageId,
		Status:      entity.AdvertStatus(dbAdvert.Status),
		CreatedAt:   dbAdvert.CreatedAt,
		UpdatedAt:   dbAdvert.UpdatedAt,
		IsSaved:     isSaved,
		IsViewed:    isViewed,
		ViewsNumber: uint(viewedCount),
		SavesNumber: uint(savedCount),
	}
}

func (r *AdvertDB) Add(a *entity.Advert) (*entity.Advert, error) {
	var dbAdvert AdvertRepoModel

	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	logger := middleware.GetLogger(r.ctx)
	logger.Info("adding advert to db", zap.Any("advert", a))

	err := r.DB.QueryRow(ctx, insertAdvertQuery,
		a.Title,
		a.Description,
		a.Price,
		a.Location,
		a.HasDelivery,
		a.CategoryId,
		a.SellerId,
		string(a.Status)).Scan(
		&dbAdvert.ID,
		&dbAdvert.Title,
		&dbAdvert.Description,
		&dbAdvert.Price,
		&dbAdvert.Location,
		&dbAdvert.HasDelivery,
		&dbAdvert.CategoryId,
		&dbAdvert.SellerId,
		&dbAdvert.ImageId,
		&dbAdvert.Status,
	)
	if err != nil {
		logger.Error("error adding advert", zap.Error(err))
		return nil, entity.PSQLWrap(err)
	}

	logger.Info("advert added", zap.Any("advert", dbAdvert))

	return &entity.Advert{
		ID:          dbAdvert.ID,
		Title:       dbAdvert.Title,
		Description: dbAdvert.Description,
		Price:       dbAdvert.Price,
		Location:    dbAdvert.Location,
		HasDelivery: dbAdvert.HasDelivery,
		CategoryId:  dbAdvert.CategoryId,
		SellerId:    dbAdvert.SellerId,
		ImageId:     dbAdvert.ImageId,
		Status:      entity.AdvertStatus(dbAdvert.Status),
	}, nil
}

func (r *AdvertDB) Get(limit, offset int, userId uuid.UUID) ([]*entity.Advert, error) {
	var adverts []*entity.Advert

	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	logger := middleware.GetLogger(r.ctx)
	logger.Info("getting adverts from db", zap.Int("limit", limit), zap.Int("offset", offset))

	rows, err := r.DB.Query(ctx, selectAdvertsQuery, limit, offset)
	if err != nil {
		logger.Error("failed to execute query", zap.Error(err))
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
			&dbAdvert.ImageId,
			&dbAdvert.Status,
			&dbAdvert.CreatedAt,
			&dbAdvert.UpdatedAt,
		); err != nil {
			logger.Error("failed to scan row", zap.Error(err))
			return nil, entity.PSQLWrap(err)
		}
		adverts = append(adverts, r.convertToEntityAdvert(dbAdvert, userId))
	}

	if err := rows.Err(); err != nil {
		logger.Error("error iterating over rows", zap.Error(err))
		return nil, entity.PSQLWrap(err)
	}

	return adverts, nil
}

func (r *AdvertDB) GetByCategoryId(categoryId, userId uuid.UUID) ([]*entity.Advert, error) {
	var adverts []*entity.Advert

	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	logger := middleware.GetLogger(r.ctx)
	logger.Info("getting adverts by category id from db", zap.String("category_id", categoryId.String()))

	rows, err := r.DB.Query(ctx, selectAdvertsByCategoryIdQuery, categoryId)
	if err != nil {
		logger.Error("failed to execute query", zap.Error(err), zap.String("category_id", categoryId.String()))
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
			&dbAdvert.ImageId,
			&dbAdvert.Status,
			&dbAdvert.CreatedAt,
			&dbAdvert.UpdatedAt,
		); err != nil {
			logger.Error("failed to scan row", zap.Error(err), zap.String("category_id", categoryId.String()))
			return nil, entity.PSQLWrap(err)
		}
		adverts = append(adverts, r.convertToEntityAdvert(dbAdvert, userId))
	}

	if err := rows.Err(); err != nil {
		logger.Error("error iterating over rows", zap.Error(err), zap.String("category_id", categoryId.String()))
		return nil, entity.PSQLWrap(err)
	}

	return adverts, nil
}

func (r *AdvertDB) GetBySellerId(sellerId, userId uuid.UUID) ([]*entity.Advert, error) {
	var adverts []*entity.Advert

	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	logger := middleware.GetLogger(r.ctx)
	logger.Info("getting adverts by seller id from db", zap.String("seller_id", sellerId.String()))

	rows, err := r.DB.Query(ctx, selectAdvertsBySellerIdQuery, sellerId)
	if err != nil {
		logger.Error("failed to execute query", zap.Error(err), zap.String("seller_id", sellerId.String()))
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
			&dbAdvert.ImageId,
			&dbAdvert.Status,
			&dbAdvert.CreatedAt,
			&dbAdvert.UpdatedAt,
		); err != nil {
			logger.Error("failed to scan row", zap.Error(err), zap.String("seller_id", sellerId.String()))
			return nil, entity.PSQLWrap(err)
		}
		adverts = append(adverts, r.convertToEntityAdvert(dbAdvert, userId))
	}

	if err := rows.Err(); err != nil {
		logger.Error("error iterating over rows", zap.Error(err), zap.String("seller_id", sellerId.String()))
		return nil, entity.PSQLWrap(err)
	}

	return adverts, nil
}

func (r *AdvertDB) GetByCartId(cartId, userId uuid.UUID) ([]*entity.Advert, error) {
	var adverts []*entity.Advert

	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	logger := middleware.GetLogger(r.ctx)
	logger.Info("getting adverts by cart id from db", zap.String("cart_id", cartId.String()))

	rows, err := r.DB.Query(ctx, selectAdvertsByCartIdQuery, cartId)
	if err != nil {
		logger.Error("failed to execute query", zap.Error(err), zap.String("cart_id", cartId.String()))
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
			&dbAdvert.ImageId,
			&dbAdvert.Status,
			&dbAdvert.CreatedAt,
			&dbAdvert.UpdatedAt,
		); err != nil {
			logger.Error("failed to scan row", zap.Error(err), zap.String("cart_id", cartId.String()))
			return nil, entity.PSQLWrap(err)
		}
		adverts = append(adverts, r.convertToEntityAdvert(dbAdvert, userId))
	}

	if err := rows.Err(); err != nil {
		logger.Error("error iterating over rows", zap.Error(err), zap.String("user_id", cartId.String()))
		return nil, entity.PSQLWrap(err)
	}

	return adverts, nil
}

func (r *AdvertDB) BeginTransaction() (pgx.Tx, error) {
	logger := middleware.GetLogger(r.ctx)
	logger.Info("beginning transaction")

	tx, err := r.DB.Begin(r.ctx)
	if err != nil {
		logger.Error("failed to begin transaction", zap.Error(err))
		return nil, err
	}
	return tx, nil
}

func (r *AdvertDB) GetById(advertId, userId uuid.UUID) (*entity.Advert, error) {
	var dbAdvert AdvertRepoModel

	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	logger := middleware.GetLogger(r.ctx)
	logger.Info("getting advert by id from db", zap.String("advert_id", advertId.String()))

	err := r.DB.QueryRow(ctx, selectAdvertByIdQuery, advertId).Scan(
		&dbAdvert.ID,
		&dbAdvert.Title,
		&dbAdvert.Description,
		&dbAdvert.Price,
		&dbAdvert.Location,
		&dbAdvert.HasDelivery,
		&dbAdvert.CategoryId,
		&dbAdvert.SellerId,
		&dbAdvert.ImageId,
		&dbAdvert.Status,
		&dbAdvert.CreatedAt,
		&dbAdvert.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Error("advert not found", zap.Error(err), zap.String("advert_id", advertId.String()))
			return nil, entity.PSQLWrap(repository.ErrAdvertNotFound)
		}
		logger.Error("failed to query advert by id", zap.Error(err), zap.String("advert_id", advertId.String()))
		return nil, entity.PSQLWrap(err)
	}

	return r.convertToEntityAdvert(dbAdvert, userId), nil
}

func (r *AdvertDB) Update(advert *entity.Advert) error {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Minute)
	defer cancel()

	logger := middleware.GetLogger(r.ctx)
	logger.Info("updating advert in db", zap.String("advert_id", advert.ID.String()))

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
		logger.Error("failed to update advert", zap.Error(err), zap.String("advert_id", advert.ID.String()))
		return entity.PSQLWrap(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		logger.Error("advert not found", zap.String("advert_id", advert.ID.String()))
		return entity.PSQLWrap(repository.ErrAdvertNotFound)
	}

	return nil
}

func (r *AdvertDB) DeleteById(advertId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	logger := middleware.GetLogger(r.ctx)
	logger.Info("deleting advert from db", zap.String("advert_id", advertId.String()))

	result, err := r.DB.Exec(ctx, deleteAdvertByIdQuery, advertId)
	if err != nil {
		logger.Error("failed to delete advert", zap.Error(err), zap.String("advert_id", advertId.String()))
		return entity.PSQLWrap(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		logger.Error("advert not found", zap.String("advert_id", advertId.String()))
		return entity.PSQLWrap(repository.ErrAdvertNotFound)
	}

	return nil
}

func (r *AdvertDB) UpdateStatus(tx pgx.Tx, advertId uuid.UUID, status entity.AdvertStatus) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	logger := middleware.GetLogger(r.ctx)
	logger.Info("updating advert status in db", zap.String("advert_id", advertId.String()))

	result, err := tx.Exec(ctx, updateAdvertStatusQuery, status, advertId)
	if err != nil {
		logger.Error("failed to update advert status", zap.Error(err), zap.String("advert_id", advertId.String()))
		return entity.PSQLWrap(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		logger.Error("advert not found", zap.String("advert_id", advertId.String()))
		return entity.PSQLWrap(repository.ErrAdvertNotFound)
	}

	return nil
}

func (r *AdvertDB) UploadImage(advertId uuid.UUID, imageId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	logger := middleware.GetLogger(r.ctx)
	logger.Info("uploading image to db", zap.String("advert_id", advertId.String()))

	result, err := r.DB.Exec(ctx, uploadImageQuery, imageId, advertId)
	if err != nil {
		logger.Error("failed to upload image", zap.Error(err), zap.String("advert_id", advertId.String()))
		return entity.PSQLWrap(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		logger.Error("advert not found", zap.String("advert_id", advertId.String()))
		return entity.PSQLWrap(repository.ErrAdvertNotFound)
	}

	return nil
}

func (r *AdvertDB) AddToSaved(advertId, userId uuid.UUID) error {
	var savedAdvert SavedAdvertRepoModel

	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	logger := middleware.GetLogger(r.ctx)
	logger.Info("adding advert to saved in db", zap.String("advert_id", advertId.String()))

	err := r.DB.QueryRow(ctx, insertSavedAdvertQuery, userId, advertId).Scan(
		&savedAdvert.ID,
		&savedAdvert.UserId,
		&savedAdvert.AdvertId,
		&savedAdvert.CreatedAt,
	)
	if err != nil {
		logger.Error("error adding advert to saved", zap.Error(err))
		return entity.PSQLWrap(err)
	}

	logger.Info("advert added to saved", zap.Any("saved_advert", savedAdvert))

	return nil
}

func (r *AdvertDB) DeleteFromSaved(userId uuid.UUID, advertId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	logger := middleware.GetLogger(r.ctx)
	logger.Info("deleting advert from saved in db", zap.String("advert_id", advertId.String()))

	result, err := r.DB.Exec(ctx, deleteSavedAdvertQuery, advertId, userId)
	if err != nil {
		logger.Error("failed to delete advert", zap.Error(err), zap.String("advert_id", advertId.String()))
		return entity.PSQLWrap(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		logger.Error("advert not found", zap.String("advert_id", advertId.String()))
		return entity.PSQLWrap(repository.ErrAdvertNotFound)
	}

	return nil
}

func (r *AdvertDB) AddViewed(userId, advertId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	logger := middleware.GetLogger(r.ctx)
	logger.Info("adding viewed advert to db", zap.String("advert_id", advertId.String()))

	var userIdToInsert interface{}
	if userId == uuid.Nil {
		userIdToInsert = nil
	} else {
		userIdToInsert = userId
	}

	_, err := r.DB.Exec(ctx, insertViewedAdvertQuery, advertId, userIdToInsert)
	if err != nil {
		logger.Error("failed to add viewed advert", zap.Error(err), zap.String("advert_id", advertId.String()))
		return entity.PSQLWrap(err)
	}

	return nil
}

func (r *AdvertDB) CheckIfExists(advertId uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	logger := middleware.GetLogger(r.ctx)
	logger.Info("checking if advert exists in db", zap.String("advert_id", advertId.String()))

	var exists bool
	err := r.DB.QueryRow(ctx, checkIfExistsQuery, advertId).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *AdvertDB) GetSavedByUserId(userId uuid.UUID) ([]*entity.Advert, error) {
	var adverts []*entity.Advert

	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	logger := middleware.GetLogger(r.ctx)
	logger.Info("getting saved adverts by user id from db", zap.String("user_id", userId.String()))

	rows, err := r.DB.Query(ctx, selectSavedAdvertsByUserIdQuery, userId)
	if err != nil {
		logger.Error("failed to execute query", zap.Error(err), zap.String("user_id", userId.String()))
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
			&dbAdvert.ImageId,
			&dbAdvert.Status,
			&dbAdvert.CreatedAt,
			&dbAdvert.UpdatedAt,
		); err != nil {
			logger.Error("failed to scan row", zap.Error(err), zap.String("user_id", userId.String()))
			return nil, entity.PSQLWrap(err)
		}
		adverts = append(adverts, r.convertToEntityAdvert(dbAdvert, userId))
	}

	if err := rows.Err(); err != nil {
		logger.Error("error iterating over rows", zap.Error(err), zap.String("user_id", userId.String()))
		return nil, entity.PSQLWrap(err)
	}

	return adverts, nil
}

func (r *AdvertDB) Search(query string, limit, offset int, userId uuid.UUID) ([]*entity.Advert, error) {
	var adverts []*entity.Advert

	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	logger := middleware.GetLogger(r.ctx)
	logger.Info("searching adverts in db", zap.String("query", query))

	rows, err := r.DB.Query(ctx, searchAdvertsQuery, query, limit, offset)
	if err != nil {
		logger.Error("failed to execute query", zap.Error(err), zap.String("query", query))
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
			&dbAdvert.ImageId,
			&dbAdvert.Status,
			&dbAdvert.CreatedAt,
			&dbAdvert.UpdatedAt,
		); err != nil {
			logger.Error("failed to scan row", zap.Error(err), zap.String("query", query))
			return nil, entity.PSQLWrap(err)
		}
		adverts = append(adverts, r.convertToEntityAdvert(dbAdvert, userId))
	}

	if err := rows.Err(); err != nil {
		logger.Error("error iterating over rows", zap.Error(err), zap.String("query", query))
		return nil, entity.PSQLWrap(err)
	}

	return adverts, nil
}

func (r *AdvertDB) Count() (int, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	logger := middleware.GetLogger(r.ctx)
	logger.Info("counting adverts in db")

	var count int
	err := r.DB.QueryRow(ctx, countAdvertsQuery).Scan(&count)
	if err != nil {
		logger.Error("failed to count adverts", zap.Error(err))
		return 0, err
	}
	return count, nil
}

func (r *AdvertDB) GetByUserId(sellerId, userId uuid.UUID) ([]*entity.Advert, error) {
	var adverts []*entity.Advert

	ctx, cancel := context.WithTimeout(r.ctx, r.timeout)
	defer cancel()

	logger := middleware.GetLogger(r.ctx)
	logger.Info("getting adverts by user id from db", zap.String("seller_id", sellerId.String()))

	rows, err := r.DB.Query(ctx, selectAdvertsByUserIdQuery, sellerId)
	if err != nil {
		logger.Error("failed to execute query", zap.Error(err), zap.String("seller_id", sellerId.String()))
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
			&dbAdvert.ImageId,
			&dbAdvert.Status,
			&dbAdvert.CreatedAt,
			&dbAdvert.UpdatedAt,
		); err != nil {
			logger.Error("failed to scan row", zap.Error(err), zap.String("seller_id", sellerId.String()))
			return nil, entity.PSQLWrap(err)
		}
		adverts = append(adverts, r.convertToEntityAdvert(dbAdvert, userId))
	}

	if err := rows.Err(); err != nil {
		logger.Error("error iterating over rows", zap.Error(err), zap.String("seller_id", sellerId.String()))
		return nil, entity.PSQLWrap(err)
	}

	return adverts, nil
}

func (r *AdvertDB) PromoteAdvert(advertID uuid.UUID) (*entity.Advert, error) {
	var advert entity.Advert
	err := r.DB.QueryRow(context.Background(), promoteAdvertQuery, advertID).Scan(
		&advert.ID,
		&advert.Title,
		&advert.PromotedUntil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to promote advert: %w", err)
	}
	return &advert, nil
}
