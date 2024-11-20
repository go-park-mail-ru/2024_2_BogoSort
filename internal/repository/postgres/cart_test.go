package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/mocks"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupCartMockDB(t *testing.T) (pgxmock.PgxPoolIface, *mocks.PgxMockAdapter) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)
	adapter := mocks.NewPgxMockAdapter(mockPool)
	return mockPool, adapter
}

func setupCartTest(t *testing.T) (pgxmock.PgxPoolIface, *mocks.PgxMockAdapter, *CartDB, func()) {
	mockPool, adapter := setupCartMockDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	repo := &CartDB{
		DB:     adapter,
		ctx:    ctx,
		logger: zap.NewNop(),
	}

	return mockPool, adapter, repo, func() {
		cancel()
		mockPool.Close()
	}
}

func TestCartDB_GetCartByUserID(t *testing.T) {
	mockPool, _, repo, teardown := setupCartTest(t)
	defer teardown()

	userID := uuid.New()
	cartID := uuid.New()

	mockPool.ExpectQuery(`SELECT id, user_id, status FROM cart WHERE user_id = \$1 AND status = 'active' LIMIT 1`).
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "status"}).
			AddRow(cartID, userID, entity.CartStatusActive))

	cart, err := repo.GetByUserId(userID)
	assert.NoError(t, err)
	assert.Equal(t, cartID, cart.ID)
	assert.Equal(t, entity.CartStatusActive, cart.Status)

	mockPool.ExpectQuery(`SELECT id, user_id, status FROM cart WHERE user_id = \$1 AND status = 'active' LIMIT 1`).
		WithArgs(userID).
		WillReturnError(pgx.ErrNoRows)

	cart, err = repo.GetByUserId(userID)
	assert.Error(t, err)
	assert.Equal(t, repository.ErrCartNotFound, err)
	assert.Equal(t, uuid.Nil, cart.ID)

	mockPool.ExpectQuery(`SELECT id, user_id, status FROM cart WHERE user_id = \$1 AND status = 'active' LIMIT 1`).
		WithArgs(userID).
		WillReturnError(errors.New("error getting cart by user id"))

	cart, err = repo.GetByUserId(userID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error getting cart by user id")
	assert.Equal(t, uuid.Nil, cart.ID)

	err = mockPool.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCartDB_DeleteAdvertFromCart(t *testing.T) {
	mockPool, _, repo, teardown := setupCartTest(t)
	defer teardown()

	cartID := uuid.New()
	advertID := uuid.New()

	mockPool.ExpectBegin()

	mockPool.ExpectExec(`DELETE FROM cart_advert WHERE cart_id = \$1 AND advert_id = \$2`).
		WithArgs(cartID, advertID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1)) // 1 строка удалена

	mockPool.ExpectCommit()

	mockTx, err := mockPool.Begin(context.Background())
	assert.NoError(t, err)

	err = repo.DeleteAdvert(cartID, advertID)
	assert.NoError(t, err)

	err = mockTx.Commit(context.Background())
	assert.NoError(t, err)

	mockPool.ExpectBegin()

	mockPool.ExpectExec(`DELETE FROM cart_advert WHERE cart_id = \$1 AND advert_id = \$2`).
		WithArgs(cartID, advertID).
		WillReturnResult(pgxmock.NewResult("DELETE", 0)) // 0 строк удалено

	mockPool.ExpectRollback()

	mockTx2, err := mockPool.Begin(context.Background())
	assert.NoError(t, err)

	err = repo.DeleteAdvert(cartID, advertID)
	assert.Error(t, err)
	assert.Equal(t, repository.ErrCartOrAdvertNotFound, err)

	err = mockTx2.Rollback(context.Background())
	assert.NoError(t, err)

	mockPool.ExpectBegin()

	mockPool.ExpectExec(`DELETE FROM cart_advert WHERE cart_id = \$1 AND advert_id = \$2`).
		WithArgs(cartID, advertID).
		WillReturnError(errors.New("delete error"))

	mockPool.ExpectRollback()

	mockTx3, err := mockPool.Begin(context.Background())
	assert.NoError(t, err)

	err = repo.DeleteAdvert(cartID, advertID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error deleting Advert from cart")

	err = mockTx3.Rollback(context.Background())
	assert.NoError(t, err)

	err = mockPool.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCartDB_GetAdvertsByCartID(t *testing.T) {
	mockPool, _, repo, teardown := setupCartTest(t)
	defer teardown()

	cartID := uuid.New()

	mockPool.ExpectQuery(`SELECT a.id, a.title, a.description, a.price, a.location, a.has_delivery, a.status FROM cart_advert ca JOIN advert a ON ca.advert_id = a.id WHERE ca.cart_id = \$1`).
		WithArgs(cartID).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "title", "description", "price", "location", "has_delivery", "status",
		}).
			AddRow(
				uuid.New(),
				"Test Advert",
				"Test Description",
				uint(100),
				"Test Location",
				true,
				entity.AdvertStatusActive,
			))

	adverts, err := repo.GetAdvertsByCartId(cartID)
	assert.NoError(t, err)
	assert.Len(t, adverts, 1)
	assert.Equal(t, "Test Advert", adverts[0].Title)

	mockPool.ExpectQuery(`SELECT a.id, a.title, a.description, a.price, a.location, a.has_delivery, a.status FROM cart_advert ca JOIN advert a ON ca.advert_id = a.id WHERE ca.cart_id = \$1`).
		WithArgs(cartID).
		WillReturnError(pgx.ErrNoRows)

	adverts, err = repo.GetAdvertsByCartId(cartID)
	assert.Error(t, err)
	assert.Equal(t, repository.ErrCartNotFound, err)
	assert.Nil(t, adverts)

	mockPool.ExpectQuery(`SELECT a.id, a.title, a.description, a.price, a.location, a.has_delivery, a.status FROM cart_advert ca JOIN advert a ON ca.advert_id = a.id WHERE ca.cart_id = \$1`).
		WithArgs(cartID).
		WillReturnError(errors.New("query error"))

	adverts, err = repo.GetAdvertsByCartId(cartID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error getting adverts by cart id")
	assert.Nil(t, adverts)

	err = mockPool.ExpectationsWereMet()
	assert.NoError(t, err)
}
