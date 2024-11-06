package postgres

import (
	"context"
	"database/sql"
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

func setupUserMockDB(t *testing.T) (pgxmock.PgxPoolIface, *mocks.PgxMockAdapter) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)
	adapter := mocks.NewPgxMockAdapter(mockPool)
	return mockPool, adapter
}

func setupUserTest(t *testing.T) (pgxmock.PgxPoolIface, *mocks.PgxMockAdapter, *UserDB, func()) {
	mockPool, adapter := setupUserMockDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	repo := &UserDB{
		DB:      adapter,
		logger:  zap.L(),
		ctx:     ctx,
		timeout: 10 * time.Second,
	}

	return mockPool, adapter, repo, func() {
		cancel()
		mockPool.Close()
	}
}

func TestUserDB_AddUser(t *testing.T) {
	mockPool, _, repo, teardown := setupUserTest(t)
	defer teardown()

	mockPool.ExpectBegin()

	tx, err := repo.BeginTransaction()
	assert.NoError(t, err)

	email := "test@example.com"
	hash := []byte("hash")
	salt := []byte("salt")

	mockPool.ExpectQuery(`INSERT INTO "user" \(email, password_hash, password_salt, status\) VALUES \(\$1, \$2, \$3, 'active'\) RETURNING id, email, password_hash, password_salt, username, phone_number, image_id, status`).
		WithArgs(email, hash, salt).
		WillReturnRows(pgxmock.NewRows([]string{"id", "email", "password_hash", "password_salt", "username", "phone_number", "image_id", "status"}).
			AddRow(uuid.New(), email, hash, salt, sql.NullString{Valid: false}, sql.NullString{Valid: false}, uuid.Nil, "active"))

	id, err := repo.AddUser(tx, email, hash, salt)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)

	mockPool.ExpectQuery(`INSERT INTO "user" \(email, password_hash, password_salt, status\)`).
		WithArgs(email, hash, salt).
		WillReturnError(errors.New("insert error"))

	_, err = repo.AddUser(tx, email, hash, salt)
	assert.Error(t, err)

	mockPool.ExpectRollback()
	assert.NoError(t, tx.Rollback(context.Background()))

	err = mockPool.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserDB_GetUserById(t *testing.T) {
	mockPool, _, repo, teardown := setupUserTest(t)
	defer teardown()

	userId := uuid.New()

	mockPool.ExpectQuery(`SELECT id, email, password_hash, password_salt, username, phone_number, image_id, status, created_at, updated_at FROM "user" WHERE id = \$1`).
		WithArgs(userId).
		WillReturnRows(pgxmock.NewRows([]string{"id", "email", "password_hash", "password_salt", "username", "phone_number", "image_id", "status", "created_at", "updated_at"}).
			AddRow(userId, "test@example.com", []byte("hash"), []byte("salt"), sql.NullString{String: "Test User", Valid: true}, sql.NullString{String: "1234567890", Valid: true}, uuid.Nil, sql.NullString{String: "active", Valid: true}, time.Now(), time.Now()))

	user, err := repo.GetUserById(userId)
	assert.NoError(t, err)
	assert.Equal(t, "Test User", user.Username)

	mockPool.ExpectQuery(`SELECT id, email, password_hash, password_salt, username, phone_number, image_id, status, created_at, updated_at FROM "user" WHERE id = \$1`).
		WithArgs(userId).
		WillReturnError(pgx.ErrNoRows)

	user, err = repo.GetUserById(userId)
	assert.Error(t, err)
	assert.Nil(t, user)

	err = mockPool.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserDB_UpdateUser(t *testing.T) {
	mockPool, _, repo, teardown := setupUserTest(t)
	defer teardown()

	user := &entity.User{
		ID:       uuid.New(),
		Username: "Updated User",
		Phone:    "9876543210",
		// другие поля...
	}

	// Ожидаем успешное обновление
	mockPool.ExpectExec(`UPDATE "user" SET username = \$1, phone_number = \$2, image_id = \$3 WHERE id = \$4`).
		WithArgs(user.Username, user.Phone, uuid.Nil, user.ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := repo.UpdateUser(user)
	assert.NoError(t, err)

	// Ожидаем, что строка не будет найдена
	mockPool.ExpectExec(`UPDATE "user" SET username = \$1, phone_number = \$2, image_id = \$3 WHERE id = \$4`).
		WithArgs(user.Username, user.Phone, uuid.Nil, user.ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err = repo.UpdateUser(user)
	assert.Error(t, err)
	assert.Equal(t, repository.ErrUserNotFound, err)

	// Ожидаем ошибку выполнения
	mockPool.ExpectExec(`UPDATE "user" SET username = \$1, phone_number = \$2, image_id = \$3 WHERE id = \$4`).
		WithArgs(user.Username, user.Phone, uuid.Nil, user.ID).
		WillReturnError(errors.New("update error"))

	err = repo.UpdateUser(user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error updating user")

	err = mockPool.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserDB_DeleteUser(t *testing.T) {
	mockPool, _, repo, teardown := setupUserTest(t)
	defer teardown()

	userId := uuid.New()

	// Ожидаем успешное удаление
	mockPool.ExpectExec(`DELETE FROM "user" WHERE id = \$1`).
		WithArgs(userId).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := repo.DeleteUser(userId)
	assert.NoError(t, err)

	// Ожидаем, что строка не будет найдена
	mockPool.ExpectExec(`DELETE FROM "user" WHERE id = \$1`).
		WithArgs(userId).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	err = repo.DeleteUser(userId)
	assert.Error(t, err)
	assert.Equal(t, repository.ErrUserNotFound, err)

	// Ожидаем ошибку выполнения
	mockPool.ExpectExec(`DELETE FROM "user" WHERE id = \$1`).
		WithArgs(userId).
		WillReturnError(errors.New("delete error"))

	err = repo.DeleteUser(userId)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error deleting user")

	err = mockPool.ExpectationsWereMet()
	assert.NoError(t, err)
}
