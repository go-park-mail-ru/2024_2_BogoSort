package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
)

const (
	queryGetUserByEmail = `
		SELECT id, email, password_hash, password_salt, username, phone_number, image_id, status
		FROM "user"
		WHERE email = $1
	`

	queryGetUserById = `
		SELECT id, email, password_hash, password_salt, username, phone_number, image_id, status, created_at, updated_at
		FROM "user"
		WHERE id = $1
	`

	queryAddUser = `
		INSERT INTO "user" (email, password_hash, password_salt, status) 
		VALUES ($1, $2, $3, 'active')
		RETURNING id, email, password_hash, password_salt, username, phone_number, image_id, status
	`

	queryUpdateUser = `
		UPDATE "user" 
		SET username = $1, phone_number = $2, image_id = $3 
		WHERE id = $4
	`

	queryDeleteUser = `
		DELETE FROM "user" WHERE id = $1
	`

	uploadAvatarQuery = `
		UPDATE "user" SET image_id = $1 WHERE id = $2
	`
)

type UsersDB struct {
	DB     *pgxpool.Pool
	ctx    context.Context
	logger *zap.Logger
}

type DBUser struct {
	ID           uuid.UUID
	Email        string
	PasswordHash []byte
	PasswordSalt []byte
	Username     sql.NullString
	Phone        sql.NullString
	AvatarId     uuid.UUID
	Status       sql.NullString
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUserRepository(db *pgxpool.Pool, ctx context.Context, logger *zap.Logger) (repository.User, error) {
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	return &UsersDB{
		DB:     db,
		ctx:    ctx,
		logger: logger,
	}, nil
}

func (us *DBUser) GetEntity() entity.User {
	return entity.User{
		ID:           us.ID,
		Email:        us.Email,
		PasswordHash: us.PasswordHash,
		PasswordSalt: us.PasswordSalt,
		Username:     us.Username.String,
		Phone:        us.Phone.String,
		AvatarId:     us.AvatarId,
		Status:       us.Status.String,
		CreatedAt:    us.CreatedAt,
		UpdatedAt:    us.UpdatedAt,
	}
}

func (us *UsersDB) BeginTransaction() (pgx.Tx, error) {
	tx, err := us.DB.Begin(us.ctx)
	if err != nil {
		us.logger.Error("failed to begin transaction", zap.Error(err))
		return nil, err
	}
	return tx, nil
}

func (us *UsersDB) AddUser(tx pgx.Tx, email string, hash, salt []byte) (uuid.UUID, error) {
	var dbUser DBUser

	err := tx.QueryRow(us.ctx, queryAddUser, email, hash, salt).Scan(
		&dbUser.ID,
		&dbUser.Email,
		&dbUser.PasswordHash,
		&dbUser.PasswordSalt,
		&dbUser.Username,
		&dbUser.Phone,
		&dbUser.AvatarId,
		&dbUser.Status,
	)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		us.logger.Error("user already exists", zap.String("email", email))
		return uuid.Nil, repository.ErrUserAlreadyExists
	case err != nil:
		us.logger.Error("error adding user", zap.String("email", email), zap.Error(err))
		return uuid.Nil, entity.PSQLWrap(errors.New("error adding user"), err)
	}

	return dbUser.ID, nil
}

func (us *UsersDB) GetUserByEmail(email string) (*entity.User, error) {
	var dbUser DBUser
	err := us.DB.QueryRow(us.ctx, queryGetUserByEmail, email).Scan(
		&dbUser.ID,
		&dbUser.Email,
		&dbUser.PasswordHash,
		&dbUser.PasswordSalt,
		&dbUser.Username,
		&dbUser.Phone,
		&dbUser.AvatarId,
		&dbUser.Status,
		&dbUser.CreatedAt,
		&dbUser.UpdatedAt,
	)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		us.logger.Error("user not found", zap.String("email", email))
		return nil, repository.ErrUserNotFound
	case err != nil:
		us.logger.Error("error getting user by email", zap.String("email", email), zap.Error(err))
		return nil, entity.PSQLWrap(errors.New("error getting user by email"), err)
	}

	user := dbUser.GetEntity()
	return &user, nil
}

func (us *UsersDB) GetUserById(id uuid.UUID) (*entity.User, error) {
	var dbUser DBUser
	err := us.DB.QueryRow(us.ctx, queryGetUserById, id).Scan(
		&dbUser.ID,
		&dbUser.Email,
		&dbUser.PasswordHash,
		&dbUser.PasswordSalt,
		&dbUser.Username,
		&dbUser.Phone,
		&dbUser.AvatarId,
		&dbUser.Status,
		&dbUser.CreatedAt,
		&dbUser.UpdatedAt,
	)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		us.logger.Error("user not found", zap.String("id", id.String()))
		return nil, repository.ErrUserNotFound
	case err != nil:
		us.logger.Error("error getting user by id", zap.String("id", id.String()), zap.Error(err))
		return nil, entity.PSQLWrap(errors.New("error getting user by id"), err)
	}

	user := dbUser.GetEntity()
	return &user, nil
}

func (us *UsersDB) UpdateUser(user *entity.User) error {
	_, err := us.DB.Exec(us.ctx, queryUpdateUser, user.Username, user.Phone, uuid.Nil, user.ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		us.logger.Error("user not found", zap.String("id", user.ID.String()))
		return repository.ErrUserNotFound
	case err != nil:
		us.logger.Error("error updating user", zap.String("id", user.ID.String()), zap.Error(err))
		return entity.PSQLWrap(errors.New("error updating user"), err)
	}

	return nil
}

func (us *UsersDB) DeleteUser(userID uuid.UUID) error {
	_, err := us.DB.Exec(us.ctx, queryDeleteUser, userID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		us.logger.Error("user not found", zap.String("id", userID.String()))
		return repository.ErrUserNotFound
	case err != nil:
		us.logger.Error("error deleting user", zap.String("id", userID.String()), zap.Error(err))
		return entity.PSQLWrap(errors.New("error deleting user"), err)
	}

	return nil
}

func (us *UsersDB) UploadImage(userID uuid.UUID, imageId uuid.UUID) error {
	result, err := us.DB.Exec(us.ctx, uploadAvatarQuery, imageId, userID)
	if err != nil {
		us.logger.Error("failed to upload image", zap.Error(err), zap.String("user_id", userID.String()))
		return entity.PSQLWrap(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		us.logger.Error("user not found", zap.String("user_id", userID.String()))
		return entity.PSQLWrap(repository.ErrUserNotFound)
	}

	return nil
}
