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
)

type UsersDB struct {
	DB *pgxpool.Pool
}

type DBUser struct {
	ID           uuid.UUID
	Email        string
	PasswordHash []byte
	PasswordSalt []byte
	Username     sql.NullString
	Phone        sql.NullString
	AvatarId     sql.NullString
	Status       sql.NullInt64
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
		Status:       uint(us.Status.Int64),
	}
}

func (us *UsersDB) GetUserByEmail(email string) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, password_salt, username, phone, avatar_id, status
		FROM users
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

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		return nil, entity.PSQLWrap(errors.New("ошибка при получении пользователя по email"), err)
	}

	user := dbUser.GetEntity()
	return &user, nil
}

func (us *UsersDB) GetUserById(id uuid.UUID) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, password_salt, username, phone, avatar_id, status
		FROM users
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
		return nil, repository.ErrUserNotFound
	case err != nil:
		return nil, entity.PSQLWrap(errors.New("ошибка при получении пользователя по id"), err)
	}

	user := dbUser.GetEntity()
	return &user, nil
}

func (us *UsersDB) AddUser(email string, hash, salt []byte) (uuid.UUID, error) {
	query := `
		INSERT INTO users (email, password_hash, password_salt) VALUES ($1, $2)
		RETURNING id, email, password_hash, password_salt, username, phone, avatar_id, status
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
		return uuid.Nil, repository.ErrUserAlreadyExists
	case err != nil:
		return uuid.Nil, entity.PSQLWrap(errors.New("ошибка при добавлении пользователя"), err)
	}

	return dbUser.ID, nil
}

func (us *UsersDB) UpdateUser(user *entity.User) error {
	query := `
		UPDATE users SET username = $1, phone = $2, avatar_id = $3, status = $4, WHERE id = $6
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := us.DB.Exec(ctx, query, user.Username, user.Phone, user.AvatarId, user.Status)
	return entity.PSQLWrap(errors.New("ошибка при обновлении пользователя"), err)
}

func (us *UsersDB) DeleteUser(userID uuid.UUID) error {
	query := `
		DELETE FROM users WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := us.DB.Exec(ctx, query, userID)
	return entity.PSQLWrap(errors.New("ошибка при удалении пользователя"), err)
}
