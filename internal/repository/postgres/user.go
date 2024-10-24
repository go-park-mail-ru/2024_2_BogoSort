package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository"
	"github.com/jackc/pgx"
)

type UsersDB struct {
	DB *pgx.Conn
}

type DBUser struct {
	ID           uint
	Email        string
	PasswordHash []byte
	Username     sql.NullString
	Phone        sql.NullString
	AvatarId     sql.NullString
	Status       sql.NullInt64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUserRepository(db *pgx.Conn) repository.User {
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

	var dbUser DBUser
	err := us.DB.QueryRow(query, email).Scan(
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

func (us *UsersDB) GetUserById(id int) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, username, phone, avatar_id, status, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var dbUser DBUser
	err := us.DB.QueryRow(query, id).Scan(
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

func (us *UsersDB) AddUser(email, password string) (*entity.User, error) {
	query := `
		INSERT INTO users (email, password_hash) VALUES ($1, $2)
		RETURNING id, email, password_hash, username, phone, avatar_id, status, created_at, updated_at
	`

	var dbUser DBUser
	err := us.DB.QueryRow(query, email, password).Scan(
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

func (us *UsersDB) UpdateUser(user *entity.User) error {
	query := `
		UPDATE users SET username = $1, phone = $2, avatar_id = $3, status = $4, updated_at = $5 WHERE id = $6
	`

	_, err := us.DB.Exec(query, user.Username, user.Phone, user.AvatarId, user.Status, user.UpdatedAt, user.ID)
	return err
}
