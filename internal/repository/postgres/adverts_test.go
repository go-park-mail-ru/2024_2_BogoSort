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

func TestUpdateAdvert(t *testing.T) {
	mockPool, _, repo, teardown := setupAdvertTest(t)
	defer teardown()

	advertID := uuid.New()
	updatedAdvert := &entity.Advert{
		ID:          advertID,
		Title:       "Updated Advert",
		Description: "Updated Description",
		Price:       150,
		Location:    "Updated Location",
		HasDelivery: false,
		CategoryId:  uuid.New(),
		SellerId:    uuid.New(),
		Status:      "inactive",
	}

	mockPool.ExpectExec(`UPDATE advert SET title = \$1, description = \$2, price = \$3, location = \$4, has_delivery = \$5, category_id = \$6, status = \$7 WHERE id = \$8`).
		WithArgs(updatedAdvert.Title, updatedAdvert.Description, updatedAdvert.Price, updatedAdvert.Location, updatedAdvert.HasDelivery, updatedAdvert.CategoryId, updatedAdvert.Status, updatedAdvert.ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := repo.Update(updatedAdvert)
	assert.NoError(t, err)

	err = mockPool.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDeleteAdvert(t *testing.T) {
	mockPool, _, repo, teardown := setupAdvertTest(t)
	defer teardown()

	advertID := uuid.New()

	mockPool.ExpectExec(`DELETE FROM advert WHERE id = \$1`).
		WithArgs(advertID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := repo.DeleteById(advertID)
	assert.NoError(t, err)

	err = mockPool.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetAdvertById(t *testing.T) {
	mockPool, _, repo, teardown := setupAdvertTest(t)
	defer teardown()

	advertID := uuid.New()

	rows := pgxmock.NewRows([]string{
		"id", "title", "description", "price", "location", "has_delivery", "category_id", "seller_id", "image_id", "status",
	}).
		AddRow(
			advertID, "Test Advert", "Test Description", 100, "Test Location", true, uuid.New(), uuid.New(), uuid.Nil, "active",
		)

	mockPool.ExpectQuery(`SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status FROM advert WHERE id = \$1`).
		WithArgs(advertID).
		WillReturnRows(rows)

	advert, err := repo.GetById(advertID, uuid.Nil)
	assert.NoError(t, err)
	assert.Equal(t, advertID, advert.ID)

	err = mockPool.ExpectationsWereMet()
	assert.NoError(t, err)
}

func setupAdvertTest(t *testing.T) (pgxmock.PgxPoolIface, *mocks.PgxMockAdapter, *AdvertDB, func()) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)
	adapter := mocks.NewPgxMockAdapter(mockPool)
	repo := &AdvertDB{
		DB:      adapter,
		ctx:     context.Background(),
		timeout: 5 * time.Second,
	}
	return mockPool, adapter, repo, func() {
		mockPool.Close()
	}
}
