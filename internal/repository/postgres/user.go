package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UsersDB struct {
	DB     *pgxpool.Pool
	logger *zap.Logger
}

type DBUser struct {
	ID           uuid.UUID
	Email        string
	PasswordHash []byte
	PasswordSalt []byte
	Username     sql.NullString
	Phone        sql.NullString
	AvatarId     sql.NullString
	Status       sql.NullString
}

func NewUserRepository(db *pgxpool.Pool) repository.User {
	return &UsersDB{
		DB: db,
	}
}

func (us *DBUser) GetEntity() entity.User {
	return entity.User{
		ID:           us.ID,
		Email:        us.Email,
		PasswordHash: us.PasswordHash,
		PasswordSalt: us.PasswordSalt,
		Username:     us.Username.String,
		Phone:        us.Phone.String,
		AvatarId:     us.AvatarId.String,
		Status:       us.Status.String,
	}
}

func (us *UsersDB) GetUserByEmail(email string) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, password_salt, username, phone_number, image_id, status
		FROM "user"
		WHERE email = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var dbUser DBUser
	err := us.DB.QueryRow(ctx, query, email).Scan(
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
		us.logger.Error("user not found", zap.String("email", email))
		return nil, repository.ErrUserNotFound
	case err != nil:
		us.logger.Error("error getting user by email", zap.String("email", email), zap.Error(err))
		return nil, entity.PSQLWrap(errors.New("error getting user by email"), err)
	}

	user := dbUser.GetEntity()
	us.logger.Info("user found", zap.String("email", email), zap.Any("user", user))
	return &user, nil
}

func (us *UsersDB) GetUserById(id uuid.UUID) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, password_salt, username, phone_number, image_id, status
		FROM "user"
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var dbUser DBUser
	err := us.DB.QueryRow(ctx, query, id).Scan(
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
		us.logger.Error("user not found", zap.String("id", id.String()))
		return nil, repository.ErrUserNotFound
	case err != nil:
		us.logger.Error("error getting user by id", zap.String("id", id.String()), zap.Error(err))
		return nil, entity.PSQLWrap(errors.New("error getting user by id"), err)
	}

	user := dbUser.GetEntity()
	us.logger.Info("user found", zap.String("id", id.String()), zap.Any("user", user))
	return &user, nil
}

func (us *UsersDB) AddUser(email string, hash, salt []byte) (uuid.UUID, error) {
	query := `
		INSERT INTO "user" (email, password_hash, password_salt, status) VALUES ($1, $2, $3, 'active')
		RETURNING id, email, password_hash, password_salt, username, phone_number, image_id, status
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var dbUser DBUser
	err := us.DB.QueryRow(ctx, query, email, hash, salt).Scan(
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

func (us *UsersDB) UpdateUser(user *entity.User) error {
	query := `
		UPDATE "user" SET username = $1, phone_number = $2, image_id = $3 WHERE id = $4
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := us.DB.Exec(ctx, query, user.Username, user.Phone, user.AvatarId, user.ID)
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
	query := `
		DELETE FROM "user" WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := us.DB.Exec(ctx, query, userID)
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
