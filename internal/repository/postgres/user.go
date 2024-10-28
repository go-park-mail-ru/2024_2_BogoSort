package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersDB struct {
	DB *pgxpool.Pool
}

type DBUser struct {
	ID           string
	Email        string
	PasswordHash []byte
	Username     sql.NullString
	Phone        sql.NullString
	AvatarId     sql.NullString
	Status       sql.NullInt64
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
		Username:     us.Username.String,
		Phone:        us.Phone.String,
		AvatarId:     us.AvatarId.String,
		Status:       uint(us.Status.Int64),
		CreatedAt:    us.CreatedAt,
		UpdatedAt:    us.UpdatedAt,
	}
}

func (us *UsersDB) GetUserByEmail(email string) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, username, phone, avatar_id, status, created_at, updated_at
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
		&dbUser.Username,
		&dbUser.Phone,
		&dbUser.AvatarId,
		&dbUser.Status,
		&dbUser.CreatedAt,
		&dbUser.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("пользователь с email %s не найден", email)
		}
		return nil, err
	}

	user := dbUser.GetEntity()
	return &user, nil
}

func (us *UsersDB) GetUserById(id string) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, username, phone, avatar_id, status, created_at, updated_at
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
		&dbUser.Username,
		&dbUser.Phone,
		&dbUser.AvatarId,
		&dbUser.Status,
		&dbUser.CreatedAt,
		&dbUser.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	user := dbUser.GetEntity()
	return &user, nil
}

func (us *UsersDB) AddUser(email, password string) (string, error) {
	query := `
		INSERT INTO users (email, password_hash) VALUES ($1, $2)
		RETURNING id, email, password_hash, username, phone, avatar_id, status, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var dbUser DBUser
	err := us.DB.QueryRow(ctx, query, email, password).Scan(
		&dbUser.ID,
		&dbUser.Email,
		&dbUser.PasswordHash,
		&dbUser.Username,
		&dbUser.Phone,
		&dbUser.AvatarId,
		&dbUser.Status,
		&dbUser.CreatedAt,
		&dbUser.UpdatedAt,
	)

	if err != nil {
		return "", err
	}

	return dbUser.ID, nil
}

func (us *UsersDB) UpdateUser(user *entity.User) error {
	query := `
		UPDATE users SET username = $1, phone = $2, avatar_id = $3, status = $4, updated_at = $5 WHERE id = $6
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := us.DB.Exec(ctx, query, user.Username, user.Phone, user.AvatarId, user.Status, user.UpdatedAt, user.ID)
	return err
}

func (us *UsersDB) DeleteUser(userID string) error {
	query := `
		DELETE FROM users WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := us.DB.Exec(ctx, query, userID)
	return err
}
