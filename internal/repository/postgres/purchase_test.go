package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/mocks"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
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

	mockPool.ExpectQuery(`INSERT INTO purchase \(seller_id, customer_id, address, status, payment_method, delivery_method, cart_id\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\) RETURNING id, seller_id, customer_id, address, status, payment_method, delivery_method, cart_id`).
		WithArgs(
			purchase.SellerID,
			purchase.CustomerID,
			purchase.Address,
			purchase.Status,
			purchase.PaymentMethod,
			purchase.DeliveryMethod,
			purchase.CartID,
		).
		WillReturnRows(pgxmock.NewRows([]string{"id", "seller_id", "customer_id", "address", "status", "payment_method", "delivery_method", "cart_id"}).
			AddRow(uuid.New(), purchase.SellerID, purchase.CustomerID, purchase.Address, purchase.Status, purchase.PaymentMethod, purchase.DeliveryMethod, purchase.CartID))

	result, err := repo.Add(tx, purchase)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assert.Equal(t, purchase.CartID, result.CartID)
	assert.Equal(t, purchase.Address, result.Address)
	assert.Equal(t, purchase.Status, result.Status)
	assert.Equal(t, purchase.PaymentMethod, result.PaymentMethod)
	assert.Equal(t, purchase.DeliveryMethod, result.DeliveryMethod)

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
