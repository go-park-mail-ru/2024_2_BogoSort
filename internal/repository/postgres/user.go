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
		SELECT id, email, password_hash, password_salt, username, phone_number, image_id, status, created_at, updated_at
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
		SET username = $1, phone_number = $2
		WHERE id = $3
	`

	queryDeleteUser = `
		DELETE FROM "user" WHERE id = $1
	`

	uploadAvatarQuery = `
		UPDATE "user" SET image_id = $1 WHERE id = $2
	`
)

type UserDB struct {
	DB      DBExecutor
	ctx     context.Context
	logger  *zap.Logger
	timeout time.Duration
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

func NewUserRepository(db *pgxpool.Pool, ctx context.Context, logger *zap.Logger, timeout time.Duration) (repository.User, error) {
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	
	return &UserDB{
		DB:      db,
		ctx:     ctx,
		logger:  logger,
		timeout: timeout,
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

func (us *UserDB) BeginTransaction() (pgx.Tx, error) {
	tx, err := us.DB.Begin(us.ctx)
	if err != nil {
		us.logger.Error("failed to begin transaction", zap.Error(err))
		return nil, err
	}
	return tx, nil
}

func (us *UserDB) Add(tx pgx.Tx, email string, hash, salt []byte) (uuid.UUID, error) {
	var dbUser DBUser

	ctx, cancel := context.WithTimeout(us.ctx, us.timeout)
	defer cancel()

	err := tx.QueryRow(ctx, queryAddUser, email, hash, salt).Scan(
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

func (us *UserDB) GetByEmail(email string) (*entity.User, error) {
	var dbUser DBUser

	ctx, cancel := context.WithTimeout(us.ctx, us.timeout)
	defer cancel()

	err := us.DB.QueryRow(ctx, queryGetUserByEmail, email).Scan(
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

func (us *UserDB) GetById(id uuid.UUID) (*entity.User, error) {
	var dbUser DBUser

	ctx, cancel := context.WithTimeout(us.ctx, us.timeout)
	defer cancel()

	err := us.DB.QueryRow(ctx, queryGetUserById, id).Scan(
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

func (us *UserDB) Update(user *entity.User) error {
	ctx, cancel := context.WithTimeout(us.ctx, us.timeout)
	defer cancel()

	ctag, err := us.DB.Exec(ctx, queryUpdateUser, user.Username, user.Phone, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			us.logger.Error("user not found", zap.String("id", user.ID.String()))
			return repository.ErrUserNotFound
		default:
			us.logger.Error("error updating user", zap.String("id", user.ID.String()), zap.Error(err))
			return entity.PSQLWrap(errors.New("error updating user"), err)
		}
	}

	if ctag.RowsAffected() == 0 {
		us.logger.Error("user not found", zap.String("id", user.ID.String()))
		return repository.ErrUserNotFound
	}

	return nil
}

func (us *UserDB) Delete(userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(us.ctx, us.timeout)
	defer cancel()

	ctag, err := us.DB.Exec(ctx, queryDeleteUser, userID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			us.logger.Error("user not found", zap.String("id", userID.String()))
			return repository.ErrUserNotFound
		default:
			us.logger.Error("error deleting user", zap.String("id", userID.String()), zap.Error(err))
			return entity.PSQLWrap(errors.New("error deleting user"), err)
		}
	}

	if ctag.RowsAffected() == 0 {
		us.logger.Error("user not found", zap.String("id", userID.String()))
		return repository.ErrUserNotFound
	}

	return nil
}

func (us *UserDB) UploadImage(userID uuid.UUID, imageId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(us.ctx, us.timeout)
	defer cancel()

	result, err := us.DB.Exec(ctx, uploadAvatarQuery, imageId, userID)
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

func (us *UserDB) CheckIfExists(userId uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(us.ctx, us.timeout)
	defer cancel()

	var exists bool
	err := us.DB.QueryRow(ctx, checkIfExistsQuery, userId).Scan(&exists)
	if err != nil {
		return false, err
	}

	return true, nil
}
