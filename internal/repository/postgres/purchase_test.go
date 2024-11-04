package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/mocks"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupPurchaseMockDB(t *testing.T) (pgxmock.PgxPoolIface, *mocks.PgxMockAdapter) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)
	adapter := mocks.NewPgxMockAdapter(mockPool)
	return mockPool, adapter
}

func setupPurchaseTest(t *testing.T) (pgxmock.PgxPoolIface, *mocks.PgxMockAdapter, *PurchaseDB, func()) {
	mockPool, adapter := setupPurchaseMockDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	repo := &PurchaseDB{
		db:      adapter,
		logger:  zap.L(),
		ctx:     ctx,
		timeout: 10 * time.Second,
	}

	return mockPool, adapter, repo, func() {
		cancel()
		mockPool.Close()
	}
}

func TestPurchaseDB_AddPurchase(t *testing.T) {
	mockPool, _, repo, teardown := setupPurchaseTest(t)
	defer teardown()

	mockPool.ExpectBegin()

	tx, err := mockPool.Begin(context.Background())
	assert.NoError(t, err)

	purchase := &entity.Purchase{
		CartID:         uuid.New(),
		Address:        "Test Address",
		Status:         "pending",
		PaymentMethod:  "credit_card",
		DeliveryMethod: "standard",
	}

	mockPool.ExpectQuery(`INSERT INTO purchase \(cart_id, adress, status, payment_method, delivery_method\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING id, cart_id, adress, status, payment_method, delivery_method`).
		WithArgs(
			purchase.CartID,
			purchase.Address,
			purchase.Status,
			purchase.PaymentMethod,
			purchase.DeliveryMethod,
		).
		WillReturnRows(pgxmock.NewRows([]string{"id", "cart_id", "adress", "status", "payment_method", "delivery_method"}).
			AddRow(uuid.New(), purchase.CartID, purchase.Address, purchase.Status, purchase.PaymentMethod, purchase.DeliveryMethod))

	result, err := repo.AddPurchase(tx, purchase)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assert.Equal(t, purchase.CartID, result.CartID)
	assert.Equal(t, purchase.Address, result.Address)
	assert.Equal(t, purchase.Status, result.Status)
	assert.Equal(t, purchase.PaymentMethod, result.PaymentMethod)
	assert.Equal(t, purchase.DeliveryMethod, result.DeliveryMethod)

	mockPool.ExpectQuery(`INSERT INTO purchase \(cart_id, adress, status, payment_method, delivery_method\)`).
		WithArgs(
			purchase.CartID,
			purchase.Address,
			purchase.Status,
			purchase.PaymentMethod,
			purchase.DeliveryMethod,
		).
		WillReturnError(errors.New("insert error"))

	_, err = repo.AddPurchase(tx, purchase)
	assert.Error(t, err)

	mockPool.ExpectRollback()
	assert.NoError(t, tx.Rollback(context.Background()))

	err = mockPool.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPurchaseDB_BeginTransaction(t *testing.T) {
	mockPool, _, repo, teardown := setupPurchaseTest(t)
	defer teardown()

	mockPool.ExpectBegin()
	tx, err := repo.BeginTransaction()
	assert.NoError(t, err)
	assert.NotNil(t, tx)

	err = mockPool.ExpectationsWereMet()
	assert.NoError(t, err)
}
