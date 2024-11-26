package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/mocks"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func setupSellerMockDB(t *testing.T) (pgxmock.PgxPoolIface, *mocks.PgxMockAdapter) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)
	adapter := mocks.NewPgxMockAdapter(mockPool)
	return mockPool, adapter
}

func setupSellerTest(t *testing.T) (pgxmock.PgxPoolIface, *mocks.PgxMockAdapter, *SellerDB, func()) {
	mockPool, adapter := setupSellerMockDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	repo := &SellerDB{
		DB:  adapter,
		ctx: ctx,
	}

	return mockPool, adapter, repo, func() {
		cancel()
		mockPool.Close()
	}
}

func TestSellerDB_AddSeller(t *testing.T) {
	mockPool, _, repo, teardown := setupSellerTest(t)
	defer teardown()

	userID := uuid.New()
	sellerID := uuid.New()

	mockPool.ExpectBegin()

	tx, err := mockPool.Begin(context.Background())
	assert.NoError(t, err)

	mockPool.ExpectQuery(`INSERT INTO "seller" \(user_id, created_at, updated_at\) VALUES \(\$1, NOW\(\), NOW\(\)\) RETURNING id, user_id, description, created_at, updated_at`).
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "description", "created_at", "updated_at"}).
			AddRow(sellerID, userID, sql.NullString{Valid: false}, time.Now(), time.Now()))

	id, err := repo.Add(tx, userID)
	assert.NoError(t, err)
	assert.Equal(t, sellerID, id)

	mockPool.ExpectCommit()
	err = tx.Commit(context.Background())
	assert.NoError(t, err)

	mockPool.ExpectBegin()

	tx2, err := mockPool.Begin(context.Background())
	assert.NoError(t, err)

	mockPool.ExpectQuery(`INSERT INTO "seller" \(user_id, created_at, updated_at\) VALUES \(\$1, NOW\(\), NOW\(\)\) RETURNING id, user_id, description, created_at, updated_at`).
		WithArgs(userID).
		WillReturnError(pgx.ErrNoRows)

	_, err = repo.Add(tx2, userID)
	assert.Error(t, err)
	assert.Equal(t, repository.ErrSellerAlreadyExists, err)

	mockPool.ExpectRollback()
	err = tx2.Rollback(context.Background())
	assert.NoError(t, err)

	mockPool.ExpectBegin()
	tx3, err := mockPool.Begin(context.Background())
	assert.NoError(t, err)

	mockPool.ExpectQuery(`INSERT INTO "seller" \(user_id, created_at, updated_at\) VALUES \(\$1, NOW\(\), NOW\(\)\) RETURNING id, user_id, description, created_at, updated_at`).
		WithArgs(userID).
		WillReturnError(errors.New("insertion error"))

	_, err = repo.Add(tx3, userID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error adding seller")

	mockPool.ExpectRollback()
	err = tx3.Rollback(context.Background())
	assert.NoError(t, err)

	err = mockPool.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSellerDB_GetSellerByID(t *testing.T) {
	mockPool, _, repo, teardown := setupSellerTest(t)
	defer teardown()

	sellerID := uuid.New()

	mockPool.ExpectQuery(`SELECT id, user_id, description, created_at, updated_at FROM "seller" WHERE id = \$1`).
		WithArgs(sellerID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "description", "created_at", "updated_at"}).
			AddRow(sellerID, uuid.New(), sql.NullString{String: "A great seller", Valid: true}, time.Now(), time.Now()))

	mockPool.ExpectQuery(`SELECT id, user_id, description, created_at, updated_at FROM "seller" WHERE id = \$1`).
		WithArgs(sellerID).
		WillReturnError(pgx.ErrNoRows)

	mockPool.ExpectQuery(`SELECT id, user_id, description, created_at, updated_at FROM "seller" WHERE id = \$1`).
		WithArgs(sellerID).
		WillReturnError(errors.New("query error"))

	seller, err := repo.GetById(sellerID)
	assert.NoError(t, err)
	assert.NotNil(t, seller)
	assert.Equal(t, "A great seller", seller.Description)

	seller, err = repo.GetById(sellerID)
	assert.Error(t, err)
	assert.Nil(t, seller)
	assert.Equal(t, repository.ErrSellerNotFound, err)

	seller, err = repo.GetById(sellerID)
	assert.Error(t, err)
	assert.Nil(t, seller)
	assert.Contains(t, err.Error(), "error getting seller by ID")

	err = mockPool.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSellerDB_GetSellerByUserID(t *testing.T) {
	mockPool, _, repo, teardown := setupSellerTest(t)
	defer teardown()

	userID := uuid.New()

	mockPool.ExpectQuery(`SELECT id, user_id, description, created_at, updated_at FROM "seller" WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "description", "created_at", "updated_at"}).
			AddRow(uuid.New(), userID, sql.NullString{String: "Another great seller", Valid: true}, time.Now(), time.Now()))

	mockPool.ExpectQuery(`SELECT id, user_id, description, created_at, updated_at FROM "seller" WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnError(pgx.ErrNoRows)

	mockPool.ExpectQuery(`SELECT id, user_id, description, created_at, updated_at FROM "seller" WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnError(errors.New("query error"))

	seller, err := repo.GetByUserId(userID)
	assert.NoError(t, err)
	assert.NotNil(t, seller)
	assert.Equal(t, "Another great seller", seller.Description)

	seller, err = repo.GetByUserId(userID)
	assert.Error(t, err)
	assert.Nil(t, seller)
	assert.Equal(t, repository.ErrSellerNotFound, err)

	seller, err = repo.GetByUserId(userID)
	assert.Error(t, err)
	assert.Nil(t, seller)
	assert.Contains(t, err.Error(), "error getting seller by user_id")

	err = mockPool.ExpectationsWereMet()
	assert.NoError(t, err)
}
