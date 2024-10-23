package postgres

import (
	"database/sql"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repo"
	"github.com/jmoiron/sqlx"
)

type UsersDB struct {
	DB *sqlx.DB
}

type DBUser struct {
	ID           uint
	Email        string
	PasswordHash []byte
	Username     sql.NullString
	Phone        string
	AvatarId     string
	Status       uint
	CreatedAt    sql.NullTime
	UpdatedAt    sql.NullTime
}

func NewUserRepository(db *sqlx.DB) repo.User {
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
		Phone:        us.Phone,
		AvatarId:     us.AvatarId,
		Status:       us.Status,
		CreatedAt:    us.CreatedAt.Time,
		UpdatedAt:    us.UpdatedAt.Time,
	}
}

func (us *UsersDB) GetUserByEmail(email string) (*entity.User, error) {

}

func (us *UsersDB) GetUserById(id int) (*entity.User, error) {

}

func (us *UsersDB) AddUser(email, password string) (*entity.User, error) {

}

func (us *UsersDB) UpdateUser(User *entity.User) error {

}
