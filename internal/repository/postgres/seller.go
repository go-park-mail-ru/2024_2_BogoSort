package postgres

import (
	"context"
	"database/sql"
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

const (
	queryAddSeller = `
		INSERT INTO "seller" (user_id, created_at, updated_at)
		VALUES ($1, NOW(), NOW())
		RETURNING id, user_id, description, created_at, updated_at
	`

	queryGetSellerByID = `
		SELECT id, user_id, description, created_at, updated_at
		FROM "seller"
		WHERE id = $1
	`

	queryGetSellerByUserID = `
		SELECT id, user_id, description, created_at, updated_at
		FROM "seller"
		WHERE user_id = $1
	`
)

type SellerDB struct {
	DB  DBExecutor
	ctx context.Context
}

type DBSeller struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Description sql.NullString
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewSellerRepository(db *pgxpool.Pool, ctx context.Context) (repository.Seller, error) {
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	return &SellerDB{
		DB:  db,
		ctx: ctx,
	}, nil
}

func (dbSeller *DBSeller) GetEntity() entity.Seller {
	return entity.Seller{
		ID:          dbSeller.ID,
		UserID:      dbSeller.UserID,
		Description: dbSeller.Description.String,
	}
}

func (s *SellerDB) Add(tx pgx.Tx, userID uuid.UUID) (uuid.UUID, error) {
	var dbSeller DBSeller
	logger := middleware.GetLogger(s.ctx)
	logger.Info("adding seller to db", zap.String("user_id", userID.String()))

	err := tx.QueryRow(s.ctx, queryAddSeller, userID).Scan(
		&dbSeller.ID,
		&dbSeller.UserID,
		&dbSeller.Description,
		&dbSeller.CreatedAt,
		&dbSeller.UpdatedAt,
	)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		logger.Error("seller already exists", zap.String("user_id", userID.String()))
		return uuid.Nil, repository.ErrSellerAlreadyExists
	case err != nil:
		logger.Error("error adding seller", zap.String("user_id", userID.String()), zap.Error(err))
		return uuid.Nil, entity.PSQLWrap(errors.New("error adding seller"), err)
	}

	return dbSeller.ID, nil
}

func (s *SellerDB) GetById(sellerID uuid.UUID) (*entity.Seller, error) {
	var dbSeller DBSeller
	logger := middleware.GetLogger(s.ctx)
	logger.Info("getting seller by id from db", zap.String("seller_id", sellerID.String()))

	err := s.DB.QueryRow(s.ctx, queryGetSellerByID, sellerID).Scan(
		&dbSeller.ID,
		&dbSeller.UserID,
		&dbSeller.Description,
		&dbSeller.CreatedAt,
		&dbSeller.UpdatedAt,
	)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		logger.Error("seller not found", zap.String("seller_id", sellerID.String()))
		return nil, repository.ErrSellerNotFound
	case err != nil:
		logger.Error("error getting seller by ID", zap.String("seller_id", sellerID.String()), zap.Error(err))
		return nil, entity.PSQLWrap(errors.New("error getting seller by ID"), err)
	}

	logger.Info("seller found", zap.String("seller_id", sellerID.String()))
	seller := dbSeller.GetEntity()
	return &seller, nil
}

func (s *SellerDB) GetByUserId(userID uuid.UUID) (*entity.Seller, error) {
	var dbSeller DBSeller
	logger := middleware.GetLogger(s.ctx)
	logger.Info("getting seller by user id from db", zap.String("user_id", userID.String()))

	err := s.DB.QueryRow(s.ctx, queryGetSellerByUserID, userID).Scan(
		&dbSeller.ID,
		&dbSeller.UserID,
		&dbSeller.Description,
		&dbSeller.CreatedAt,
		&dbSeller.UpdatedAt,
	)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		logger.Error("seller not found by user_id", zap.String("user_id", userID.String()))
		return nil, repository.ErrSellerNotFound
	case err != nil:
		logger.Error("error getting seller by user_id", zap.String("user_id", userID.String()), zap.Error(err))
		return nil, entity.PSQLWrap(errors.New("error getting seller by user_id"), err)
	}

	seller := dbSeller.GetEntity()
	return &seller, nil
}
