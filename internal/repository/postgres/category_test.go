package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/mocks"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func setupCategoryMockDB(t *testing.T) (pgxmock.PgxPoolIface, *mocks.PgxMockAdapter) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)
	adapter := mocks.NewPgxMockAdapter(mockPool)
	return mockPool, adapter
}

func setupCategoryTest(t *testing.T) (pgxmock.PgxPoolIface, *mocks.PgxMockAdapter, *CategoryDB, func()) {
	mockPool, adapter := setupCategoryMockDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	repo := &CategoryDB{
		DB:      adapter,
		ctx:     ctx,
		timeout: 10 * time.Second,
	}

	return mockPool, adapter, repo, func() {
		cancel()
		mockPool.Close()
	}
}

func TestCategoryDB_GetCategories(t *testing.T) {
	mockPool, _, repo, teardown := setupCategoryTest(t)
	defer teardown()

	// Successful case
	mockPool.ExpectQuery(getCategoryQuery).
		WillReturnRows(pgxmock.NewRows([]string{"id", "title"}).
			AddRow(uuid.New(), "Test Category"))

	categories, err := repo.Get()
	assert.NoError(t, err)
	assert.Len(t, categories, 1)

	// Error case
	mockPool.ExpectQuery(getCategoryQuery).
		WillReturnError(errors.New("query error"))

	categories, err = repo.Get()
	assert.Error(t, err)
	assert.Nil(t, categories)

	err = mockPool.ExpectationsWereMet()
	assert.NoError(t, err)
}
