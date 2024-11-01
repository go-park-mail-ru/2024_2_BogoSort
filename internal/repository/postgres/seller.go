package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

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
	DB     *pgxpool.Pool
	ctx    context.Context
	logger *zap.Logger
}

type DBSeller struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Description sql.NullString
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewSellerRepository(db *pgxpool.Pool, ctx context.Context, logger *zap.Logger) repository.Seller {
	return &SellerDB{
		DB:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (s *SellerDB) AddSeller(tx pgx.Tx, userID uuid.UUID) (uuid.UUID, error) {
	var dbSeller DBSeller
	err := tx.QueryRow(s.ctx, queryAddSeller, userID).Scan(
		&dbSeller.ID,
		&dbSeller.UserID,
		&dbSeller.Description,
		&dbSeller.CreatedAt,
		&dbSeller.UpdatedAt,
	)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		s.logger.Error("seller already exists", zap.String("user_id", userID.String()))
		return uuid.Nil, repository.ErrSellerAlreadyExists
	case err != nil:
		s.logger.Error("error adding seller", zap.String("user_id", userID.String()), zap.Error(err))
		return uuid.Nil, entity.PSQLWrap(errors.New("error adding seller"), err)
	}

	return dbSeller.ID, nil
}

func (s *SellerDB) GetSellerByID(sellerID uuid.UUID) (*entity.Seller, error) {
	var dbSeller DBSeller
	err := s.DB.QueryRow(s.ctx, queryGetSellerByID, sellerID).Scan(
		&dbSeller.ID,
		&dbSeller.UserID,
		&dbSeller.Description,
		&dbSeller.CreatedAt,
		&dbSeller.UpdatedAt,
	)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		s.logger.Error("seller not found", zap.String("seller_id", sellerID.String()))
		return nil, repository.ErrSellerNotFound
	case err != nil:
		s.logger.Error("error getting seller by ID", zap.String("seller_id", sellerID.String()), zap.Error(err))
		return nil, err
	}

	seller := dbSeller.GetEntity()
	return &seller, nil
}

func (s *SellerDB) GetSellerByUserID(userID uuid.UUID) (*entity.Seller, error) {
	var dbSeller DBSeller
	err := s.DB.QueryRow(s.ctx, queryGetSellerByUserID, userID).Scan(
		&dbSeller.ID,
		&dbSeller.UserID,
		&dbSeller.Description,
		&dbSeller.CreatedAt,
		&dbSeller.UpdatedAt,
	)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		s.logger.Error("seller not found by user_id", zap.String("user_id", userID.String()))
		return nil, repository.ErrSellerNotFound
	case err != nil:
		s.logger.Error("error getting seller by user_id", zap.String("user_id", userID.String()), zap.Error(err))
		return nil, err
	}

	seller := dbSeller.GetEntity()
	return &seller, nil
}

func (dbSeller *DBSeller) GetEntity() entity.Seller {
	return entity.Seller{
		ID:          dbSeller.ID,
		UserID:      dbSeller.UserID,
		Description: dbSeller.Description.String,
	}
}
