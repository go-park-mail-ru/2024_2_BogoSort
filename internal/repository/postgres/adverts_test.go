package postgres

import (
	"context"
	"regexp"
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
		"id", "title", "description", "price", "location", "has_delivery", "category_id", "seller_id", "image_id", "status", "created_at", "updated_at",
	}).AddRow(
		advertID, "Test Advert", "Test Description", uint(100), "Test Location", true, uuid.New(), uuid.New(), uuid.Nil, "active", time.Now(), time.Now(),
	)

	mockPool.ExpectQuery(`SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status, created_at, updated_at FROM advert WHERE id = \$1`).
		WithArgs(advertID).
		WillReturnRows(rows)

	advert, err := repo.GetById(advertID, uuid.Nil)

	assert.NoError(t, err)
	assert.Nil(t, advert)
}

func TestAddAdvert(t *testing.T) {
	mockPool, _, repo, teardown := setupAdvertTest(t)
	defer teardown()

	newAdvert := &entity.Advert{
		Title:       "New Advert",
		Description: "New Description",
		Price:       200,
		Location:    "New Location",
		HasDelivery: true,
		CategoryId:  uuid.New(),
		SellerId:    uuid.New(),
		Status:      entity.AdvertStatusActive,
	}

	mockPool.ExpectQuery(`INSERT INTO advert \(title, description, price, location, has_delivery, category_id, seller_id, status\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\) RETURNING id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status`).
		WithArgs(newAdvert.Title, newAdvert.Description, newAdvert.Price, newAdvert.Location, newAdvert.HasDelivery, newAdvert.CategoryId, newAdvert.SellerId, "active").
		WillReturnRows(pgxmock.NewRows([]string{"id", "title", "description", "price", "location", "has_delivery", "category_id", "seller_id", "image_id", "status"}).AddRow(uuid.New(), newAdvert.Title, newAdvert.Description, newAdvert.Price, newAdvert.Location, newAdvert.HasDelivery, newAdvert.CategoryId, newAdvert.SellerId, uuid.Nil, "active"))

	result, err := repo.Add(newAdvert)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newAdvert.Title, result.Title)

	err = mockPool.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateAdvertStatus(t *testing.T) {
	mockPool, _, repo, teardown := setupAdvertTest(t)
	defer teardown()

	advertID := uuid.New()
	newStatus := entity.AdvertStatusInactive // Измените тип статуса на entity.AdvertStatus
	mockPool.ExpectBegin()
	tx, err := repo.DB.Begin(context.Background()) // Начало транзакции
	assert.NoError(t, err)

	mockPool.ExpectExec(`UPDATE advert SET status = \$1 WHERE id = \$2`).
		WithArgs(newStatus, advertID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.UpdateStatus(tx, advertID, newStatus) // Добавьте tx как аргумент
	assert.NoError(t, err)

	err = mockPool.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGet(t *testing.T) {
	mockPool, _, repo, teardown := setupAdvertTest(t)
	defer teardown()

	userID := uuid.New()
	limit, offset := 10, 0

	rows := pgxmock.NewRows([]string{"id", "title", "description", "price", "location", "has_delivery", "category_id", "seller_id", "image_id", "status", "created_at", "updated_at"}).
		AddRow(uuid.New(), "Advert 1", "Description 1", uint(100), "Location 1", true, uuid.New(), uuid.New(), uuid.Nil, "active", time.Now(), time.Now())

	mockPool.ExpectQuery(regexp.QuoteMeta(selectAdvertsQuery)).
		WithArgs(limit, offset).
		WillReturnRows(rows)

	adverts, err := repo.Get(limit, offset, userID)
	assert.NoError(t, err)
	assert.Len(t, adverts, 1)
}

func TestGetByCategoryId(t *testing.T) {
	mockPool, _, repo, teardown := setupAdvertTest(t)
	defer teardown()

	categoryID, userID := uuid.New(), uuid.New()

	rows := pgxmock.NewRows([]string{"id", "title", "description", "price", "location", "has_delivery", "category_id", "seller_id", "image_id", "status", "created_at", "updated_at"}).
		AddRow(uuid.New(), "Advert 1", "Description 1", uint(100), "Location 1", true, categoryID, uuid.New(), uuid.Nil, "active", time.Now(), time.Now())

	mockPool.ExpectQuery(regexp.QuoteMeta(selectAdvertsByCategoryIdQuery)).
		WithArgs(categoryID).
		WillReturnRows(rows)

	adverts, err := repo.GetByCategoryId(categoryID, userID)
	assert.NoError(t, err)
	assert.Len(t, adverts, 1)
}

func TestGetBySellerId(t *testing.T) {
	mockPool, _, repo, teardown := setupAdvertTest(t)
	defer teardown()

	sellerID, userID := uuid.New(), uuid.New()

	rows := pgxmock.NewRows([]string{"id", "title", "description", "price", "location", "has_delivery", "category_id", "seller_id", "image_id", "status", "created_at", "updated_at"}).
		AddRow(uuid.New(), "Advert 1", "Description 1", uint(100), "Location 1", true, uuid.New(), sellerID, uuid.Nil, "active", time.Now(), time.Now())

	mockPool.ExpectQuery(regexp.QuoteMeta(selectAdvertsBySellerIdQuery)).
		WithArgs(sellerID).
		WillReturnRows(rows)

	adverts, err := repo.GetBySellerId(sellerID, userID)
	assert.NoError(t, err)
	assert.Len(t, adverts, 1)
}

func TestSearch(t *testing.T) {
	mockPool, _, repo, teardown := setupAdvertTest(t)
	defer teardown()

	query := "test"
	limit, offset := 10, 0
	userID := uuid.New()

	rows := pgxmock.NewRows([]string{"id", "title", "description", "price", "location", "has_delivery", "category_id", "seller_id", "image_id", "status", "created_at", "updated_at"}).
		AddRow(uuid.New(), "Advert 1", "Description 1", uint(100), "Location 1", true, uuid.New(), uuid.New(), uuid.Nil, "active", time.Now(), time.Now())

	mockPool.ExpectQuery(regexp.QuoteMeta(searchAdvertsQuery)).
		WithArgs(query, limit, offset).
		WillReturnRows(rows)

	adverts, err := repo.Search(query, limit, offset, userID)
	assert.NoError(t, err)
	assert.Len(t, adverts, 1)
}

func TestGetByCartId(t *testing.T) {
	mockPool, _, repo, teardown := setupAdvertTest(t)
	defer teardown()

	cartId, userId := uuid.New(), uuid.New()

	rows := pgxmock.NewRows([]string{"id", "title", "description", "price", "location", "has_delivery", "category_id", "seller_id", "image_id", "status", "created_at", "updated_at"}).
		AddRow(uuid.New(), "Advert 1", "Description 1", uint(100), "Location 1", true, uuid.New(), uuid.New(), uuid.Nil, "active", time.Now(), time.Now())

	mockPool.ExpectQuery(regexp.QuoteMeta(selectAdvertsByCartIdQuery)).
		WithArgs(cartId).
		WillReturnRows(rows)

	adverts, err := repo.GetByCartId(cartId, userId)
	assert.NoError(t, err)
	assert.Len(t, adverts, 1)
}

func TestGetSavedByUserId(t *testing.T) {
	mockPool, _, repo, teardown := setupAdvertTest(t)
	defer teardown()

	userId := uuid.New()

	rows := pgxmock.NewRows([]string{"id", "title", "description", "price", "location", "has_delivery", "category_id", "seller_id", "image_id", "status", "created_at", "updated_at"}).
		AddRow(uuid.New(), "Advert 1", "Description 1", uint(100), "Location 1", true, uuid.New(), uuid.New(), uuid.Nil, "active", time.Now(), time.Now())

	mockPool.ExpectQuery(regexp.QuoteMeta(selectSavedAdvertsByUserIdQuery)).
		WithArgs(userId).
		WillReturnRows(rows)

	adverts, err := repo.GetSavedByUserId(userId)
	assert.NoError(t, err)
	assert.Len(t, adverts, 1)
}

func TestDeleteFromSaved(t *testing.T) {
	mockPool, _, repo, teardown := setupAdvertTest(t)
	defer teardown()

	userId, advertId := uuid.New(), uuid.New()

	mockPool.ExpectExec(regexp.QuoteMeta(deleteSavedAdvertQuery)).
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := repo.DeleteFromSaved(userId, advertId)
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
