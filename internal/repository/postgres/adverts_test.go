package postgres

// import (
// 	"context"
// 	"errors"
// 	"testing"
// 	"time"

// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
// 	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/mocks"
// 	"github.com/google/uuid"
// 	"github.com/jackc/pgx/v5"
// 	"github.com/pashagolub/pgxmock/v4"
// 	"github.com/stretchr/testify/assert"
// 	"go.uber.org/zap"
// )

// func setupAdvertMockDB(t *testing.T) (pgxmock.PgxPoolIface, *mocks.PgxMockAdapter) {
// 	mockPool, err := pgxmock.NewPool()
// 	assert.NoError(t, err)
// 	adapter := mocks.NewPgxMockAdapter(mockPool)
// 	return mockPool, adapter
// }

// func setupAdvertTest(t *testing.T) (pgxmock.PgxPoolIface, *mocks.PgxMockAdapter, *AdvertDB, func()) {
// 	mockPool, adapter := setupAdvertMockDB(t)
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

// 	repo := &AdvertDB{
// 		DB:      adapter,
// 		logger:  zap.L(),
// 		ctx:     ctx,
// 		timeout: 10 * time.Second,
// 	}

// 	return mockPool, adapter, repo, func() {
// 		cancel()
// 		mockPool.Close()
// 	}
// }

// func TestAdvertDB_AddAdvert(t *testing.T) {
// 	mockPool, _, repo, teardown := setupAdvertTest(t)
// 	defer teardown()

// 	advert := &entity.Advert{
// 		Title:       "Test Advert",
// 		Description: "Test Description",
// 		Price:       uint(100),
// 		Location:    "Test Location",
// 		HasDelivery: true,
// 		CategoryId:  uuid.New(),
// 		SellerId:    uuid.New(),
// 		ImageURL:    uuid.NullUUID{UUID: uuid.New(), Valid: true},
// 		Status:      entity.AdvertStatus("active"),
// 	}

// 	stringStatus := string(advert.Status)

// 	mockPool.ExpectQuery(`INSERT INTO advert \(title, description, price, location, has_delivery, category_id, seller_id, image_id, status\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9\) RETURNING id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status`).
// 		WithArgs(
// 			advert.Title,
// 			advert.Description,
// 			advert.Price,
// 			advert.Location,
// 			advert.HasDelivery,
// 			advert.CategoryId,
// 			advert.SellerId,
// 			advert.ImageURL,
// 			stringStatus,
// 		).
// 		WillReturnRows(pgxmock.NewRows([]string{"id", "title", "description", "price", "location", "has_delivery", "category_id", "seller_id", "image_id", "status"}).
// 			AddRow(
// 				uuid.New(),
// 				advert.Title,
// 				advert.Description,
// 				advert.Price,
// 				advert.Location,
// 				advert.HasDelivery,
// 				advert.CategoryId,
// 				advert.SellerId,
// 				advert.ImageURL,
// 				stringStatus,
// 			))

// 	result, err := repo.AddAdvert(advert)
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}

// 	assert.Equal(t, advert.Title, result.Title)
// 	assert.Equal(t, advert.Description, result.Description)
// 	assert.Equal(t, advert.Price, result.Price)
// 	assert.Equal(t, advert.Location, result.Location)
// 	assert.Equal(t, advert.HasDelivery, result.HasDelivery)
// 	assert.Equal(t, advert.CategoryId, result.CategoryId)
// 	assert.Equal(t, advert.SellerId, result.SellerId)
// 	assert.Equal(t, advert.ImageURL, result.ImageURL)
// 	assert.Equal(t, advert.Status, result.Status)

// 	mockPool.ExpectQuery(`INSERT INTO advert \(title, description, price, location, has_delivery, category_id, seller_id, image_id, status\)`).
// 		WithArgs(
// 			advert.Title,
// 			advert.Description,
// 			advert.Price,
// 			advert.Location,
// 			advert.HasDelivery,
// 			advert.CategoryId,
// 			advert.SellerId,
// 			advert.ImageURL,
// 			stringStatus,
// 		).
// 		WillReturnError(errors.New("insert error"))

// 	_, err = repo.AddAdvert(advert)
// 	assert.Error(t, err)

// 	err = mockPool.ExpectationsWereMet()
// 	assert.NoError(t, err)
// }

// func TestAdvertDB_GetAdverts(t *testing.T) {
// 	mockPool, _, repo, teardown := setupAdvertTest(t)
// 	defer teardown()

// 	// Successful case
// 	mockPool.ExpectQuery("SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status, created_at, updated_at").
// 		WithArgs(10, 0).
// 		WillReturnRows(pgxmock.NewRows([]string{"id", "title", "description", "price", "location", "has_delivery", "category_id", "seller_id", "image_id", "status", "created_at", "updated_at"}).
// 			AddRow(uuid.New(), "Test Advert", "Test Description", uint(100), "Test Location", true, uuid.New(), uuid.New(), uuid.NullUUID{}, "active", time.Now(), time.Now()))

// 	adverts, err := repo.GetAdverts(10, 0)
// 	assert.NoError(t, err)
// 	assert.Len(t, adverts, 1)

// 	// Error case
// 	mockPool.ExpectQuery("SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status, created_at, updated_at").
// 		WithArgs(10, 0).
// 		WillReturnError(errors.New("query error"))

// 	adverts, err = repo.GetAdverts(10, 0)
// 	assert.Error(t, err)
// 	assert.Nil(t, adverts)

// 	err = mockPool.ExpectationsWereMet()
// 	assert.NoError(t, err)
// }

// func TestAdvertDB_GetAdvertById(t *testing.T) {
// 	mockPool, _, repo, teardown := setupAdvertTest(t)
// 	defer teardown()

// 	advertId := uuid.New()

// 	// Successful case
// 	mockPool.ExpectQuery("SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status, created_at, updated_at").
// 		WithArgs(advertId).
// 		WillReturnRows(pgxmock.NewRows([]string{"id", "title", "description", "price", "location", "has_delivery", "category_id", "seller_id", "image_id", "status", "created_at", "updated_at"}).
// 			AddRow(advertId, "Test Advert", "Test Description", uint(100), "Test Location", true, uuid.New(), uuid.New(), uuid.NullUUID{}, "active", time.Now(), time.Now()))

// 	advert, err := repo.GetAdvertById(advertId)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "Test Advert", advert.Title)

// 	// Error case (not found)
// 	mockPool.ExpectQuery("SELECT id, title, description, price").
// 		WithArgs(advertId).
// 		WillReturnError(pgx.ErrNoRows)

// 	advert, err = repo.GetAdvertById(advertId)
// 	assert.Error(t, err)
// 	assert.Nil(t, advert)

// 	err = mockPool.ExpectationsWereMet()
// 	assert.NoError(t, err)
// }

// func TestAdvertDB_UpdateAdvert(t *testing.T) {
// 	mockPool, _, repo, teardown := setupAdvertTest(t)
// 	defer teardown()

// 	advert := &entity.Advert{
// 		ID:          uuid.New(),
// 		Title:       "Updated Advert",
// 		Description: "Updated Description",
// 		Price:       150,
// 		Location:    "Updated Location",
// 		HasDelivery: false,
// 		CategoryId:  uuid.New(),
// 		SellerId:    uuid.New(),
// 		ImageURL:    uuid.NullUUID{},
// 		Status:      entity.AdvertStatus("inactive"),
// 	}

// 	// Successful case
// 	mockPool.ExpectExec("UPDATE advert").
// 		WithArgs(advert.Title, advert.Description, advert.Price, advert.Location, advert.HasDelivery, advert.CategoryId, advert.Status, advert.ID).
// 		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

// 	err := repo.UpdateAdvert(advert)
// 	assert.NoError(t, err)

// 	// Error case (not found)
// 	mockPool.ExpectExec("UPDATE advert").
// 		WithArgs(advert.Title, advert.Description, advert.Price, advert.Location, advert.HasDelivery, advert.CategoryId, advert.Status, advert.ID).
// 		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

// 	err = repo.UpdateAdvert(advert)
// 	assert.Error(t, err)

// 	// Error case (SQL error)
// 	mockPool.ExpectExec("UPDATE advert").
// 		WithArgs(advert.Title, advert.Description, advert.Price, advert.Location, advert.HasDelivery, advert.CategoryId, advert.Status, advert.ID).
// 		WillReturnError(errors.New("update error"))

// 	err = repo.UpdateAdvert(advert)
// 	assert.Error(t, err)

// 	err = mockPool.ExpectationsWereMet()
// 	assert.NoError(t, err)
// }

// func TestAdvertDB_DeleteAdvertById(t *testing.T) {
// 	mockPool, _, repo, teardown := setupAdvertTest(t)
// 	defer teardown()

// 	advertId := uuid.New()

// 	// Successful case
// 	mockPool.ExpectExec("DELETE FROM advert").
// 		WithArgs(advertId).
// 		WillReturnResult(pgxmock.NewResult("DELETE", 1))

// 	err := repo.DeleteAdvertById(advertId)
// 	assert.NoError(t, err)

// 	// Error case (not found)
// 	mockPool.ExpectExec("DELETE FROM advert").
// 		WithArgs(advertId).
// 		WillReturnResult(pgxmock.NewResult("DELETE", 0))

// 	err = repo.DeleteAdvertById(advertId)
// 	assert.Error(t, err)

// 	// Error case (SQL error)
// 	mockPool.ExpectExec("DELETE FROM advert").
// 		WithArgs(advertId).
// 		WillReturnError(errors.New("delete error"))

// 	err = repo.DeleteAdvertById(advertId)
// 	assert.Error(t, err)

// 	err = mockPool.ExpectationsWereMet()
// 	assert.NoError(t, err)
// }

// // func TestAdvertDB_UpdateAdvertStatus(t *testing.T) {
// // 	mockPool, _, repo, teardown := setupAdvertTest(t)
// // 	defer teardown()

// // 	advertId := uuid.New()

// // 	// Successful case
// // 	mockPool.ExpectExec(`UPDATE advert SET status = \$1 WHERE id = \$2`).
// // 		WithArgs("inactive", advertId).
// // 		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

// // 	err := repo.UpdateAdvertStatus(advertId, "inactive")
// // 	assert.NoError(t, err)

// // 	// Error case (not found)
// // 	mockPool.ExpectExec(`UPDATE advert SET status = \$1 WHERE id = \$2`).
// // 		WithArgs("inactive", advertId).
// // 		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

// // 	err = repo.UpdateAdvertStatus(advertId, "inactive")
// // 	assert.Error(t, err)

// // 	// Error case (SQL error)
// // 	mockPool.ExpectExec(`UPDATE advert SET status = \$1 WHERE id = \$2`).
// // 		WithArgs("inactive", advertId).
// // 		WillReturnError(errors.New("update status error"))

// // 	err = repo.UpdateAdvertStatus(advertId, "inactive")
// // 	assert.Error(t, err)

// // 	err = mockPool.ExpectationsWereMet()
// // 	assert.NoError(t, err)
// // }

// func TestAdvertDB_UploadImage(t *testing.T) {
// 	mockPool, _, repo, teardown := setupAdvertTest(t)
// 	defer teardown()

// 	advertId := uuid.New()
// 	imageId := uuid.New()

// 	// Successful case
// 	mockPool.ExpectExec(`UPDATE advert SET image_id = \$1 WHERE id = \$2`).
// 		WithArgs(imageId, advertId).
// 		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

// 	err := repo.UploadImage(advertId, imageId)
// 	assert.NoError(t, err)

// 	// Error case (not found)
// 	mockPool.ExpectExec(`UPDATE advert SET image_id = \$1 WHERE id = \$2`).
// 		WithArgs(imageId, advertId).
// 		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

// 	err = repo.UploadImage(advertId, imageId)
// 	assert.Error(t, err)

// 	// Error case (SQL error)
// 	mockPool.ExpectExec(`UPDATE advert SET image_id = \$1 WHERE id = \$2`).
// 		WithArgs(imageId, advertId).
// 		WillReturnError(errors.New("upload image error"))

// 	err = repo.UploadImage(advertId, imageId)
// 	assert.Error(t, err)

// 	err = mockPool.ExpectationsWereMet()
// 	assert.NoError(t, err)
// }

// func TestAdvertDB_GetAdvertsByCategoryId(t *testing.T) {
// 	mockPool, _, repo, teardown := setupAdvertTest(t)
// 	defer teardown()

// 	categoryId := uuid.New()

// 	// Successful case
// 	mockPool.ExpectQuery("SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status, created_at, updated_at").
// 		WithArgs(categoryId).
// 		WillReturnRows(pgxmock.NewRows([]string{"id", "title", "description", "price", "location", "has_delivery", "category_id", "seller_id", "image_id", "status", "created_at", "updated_at"}).
// 			AddRow(uuid.New(), "Test Advert", "Test Description", uint(100), "Test Location", true, categoryId, uuid.New(), uuid.NullUUID{}, "active", time.Now(), time.Now()))

// 	adverts, err := repo.GetAdvertsByCategoryId(categoryId)
// 	assert.NoError(t, err)
// 	assert.Len(t, adverts, 1)

// 	// Error case
// 	mockPool.ExpectQuery("SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status, created_at, updated_at").
// 		WithArgs(categoryId).
// 		WillReturnError(errors.New("query error"))

// 	adverts, err = repo.GetAdvertsByCategoryId(categoryId)
// 	assert.Error(t, err)
// 	assert.Nil(t, adverts)

// 	err = mockPool.ExpectationsWereMet()
// 	assert.NoError(t, err)
// }

// func TestAdvertDB_GetAdvertsBySellerId(t *testing.T) {
// 	mockPool, _, repo, teardown := setupAdvertTest(t)
// 	defer teardown()

// 	sellerId := uuid.New()

// 	// Successful case
// 	mockPool.ExpectQuery("SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status, created_at, updated_at").
// 		WithArgs(sellerId).
// 		WillReturnRows(pgxmock.NewRows([]string{"id", "title", "description", "price", "location", "has_delivery", "category_id", "seller_id", "image_id", "status", "created_at", "updated_at"}).
// 			AddRow(uuid.New(), "Test Advert", "Test Description", uint(100), "Test Location", true, uuid.New(), sellerId, uuid.NullUUID{}, "active", time.Now(), time.Now()))

// 	adverts, err := repo.GetAdvertsBySellerId(sellerId)
// 	assert.NoError(t, err)
// 	assert.Len(t, adverts, 1)

// 	// Error case
// 	mockPool.ExpectQuery("SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status, created_at, updated_at").
// 		WithArgs(sellerId).
// 		WillReturnError(errors.New("query error"))

// 	adverts, err = repo.GetAdvertsBySellerId(sellerId)
// 	assert.Error(t, err)
// 	assert.Nil(t, adverts)

// 	err = mockPool.ExpectationsWereMet()
// 	assert.NoError(t, err)
// }

// func TestAdvertDB_GetAdvertsByCartId(t *testing.T) {
// 	mockPool, _, repo, teardown := setupAdvertTest(t)
// 	defer teardown()

// 	cartId := uuid.New()

// 	// Successful case
// 	mockPool.ExpectQuery("SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status, created_at, updated_at").
// 		WithArgs(cartId).
// 		WillReturnRows(pgxmock.NewRows([]string{"id", "title", "description", "price", "location", "has_delivery", "category_id", "seller_id", "image_id", "status", "created_at", "updated_at"}).
// 			AddRow(uuid.New(), "Test Advert", "Test Description", uint(100), "Test Location", true, uuid.New(), uuid.New(), uuid.NullUUID{}, "active", time.Now(), time.Now()))

// 	adverts, err := repo.GetAdvertsByCartId(cartId)
// 	assert.NoError(t, err)
// 	assert.Len(t, adverts, 1)

// 	// Error case
// 	mockPool.ExpectQuery("SELECT id, title, description, price, location, has_delivery, category_id, seller_id, image_id, status, created_at, updated_at").
// 		WithArgs(cartId).
// 		WillReturnError(errors.New("query error"))

// 	adverts, err = repo.GetAdvertsByCartId(cartId)
// 	assert.Error(t, err)
// 	assert.Nil(t, adverts)

// 	err = mockPool.ExpectationsWereMet()
// 	assert.NoError(t, err)
// }
