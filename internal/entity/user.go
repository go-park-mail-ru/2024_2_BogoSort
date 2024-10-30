package entity

import (
	"bytes"
	"errors"
	"regexp"
	_ "time"
	"unicode/utf8"

	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/random"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

const (
	PasswordHashTime      = 1
	PasswordHashKibMemory = 64 * 1024
	PasswordHashThreads   = 4
)

type User struct {
	ID           uuid.UUID `db:"uuid"`
	Email        string    `db:"email"`
	PasswordHash []byte    `db:"password_hash"`
	PasswordSalt []byte    `db:"password_salt"`
	Username     string    `db:"username"`
	Phone        string    `db:"phone"`
	AvatarId     string    `db:"avatar_id"`
	Status       string    `db:"status" default:"active"`
}

func ValidatePassword(password string) error {
	switch {
	case utf8.RuneCountInString(password) < 8:
		return errors.New("password must contain at least 8 characters")
	case utf8.RuneCountInString(password) > 32:
		return errors.New("password cannot be longer than 32 characters")
	case !regexp.MustCompile(`^[!@#$%^&*()_+\-=.,\w]+$`).MatchString(password):
		return errors.New("password can consist of latin letters, numbers and " +
			"special characters !@#$%^&*()_+\\-=")
	default:
		return nil
	}
}

func ValidateEmail(email string) error {
	re := regexp.MustCompile("^([a-z0-9!#$%&'*+\\\\/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+\\\\/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)$") // nolint: lll

	if !re.MatchString(email) {
		return errors.New("invalid email")
	}
	if len(email) > 256 {
		return errors.New("email cannot be longer than 256 characters")
	}
	return nil
}

func HashPassword(password string) (salt []byte, hash []byte, err error) {
	salt, err = random.Bytes(8)
	if err != nil {
		return nil, nil, errors.New("unexpected error")
	}
	hash = argon2.IDKey(
		[]byte(password),
		salt,
		PasswordHashTime,
		PasswordHashKibMemory,
		PasswordHashThreads,
		32,
	)
	return salt, hash, nil
}

func (u *User) CheckPassword(password string) bool {
	return bytes.Equal(
		argon2.IDKey(
			[]byte(password),
			u.PasswordSalt,
			PasswordHashTime,
			PasswordHashKibMemory,
			PasswordHashThreads,
			32,
		),
		u.PasswordHash,
	)
}

func ValidateName(name string) error {
	if utf8.RuneCountInString(name) > 30 {
		return errors.New("name cannot be longer than 30 characters")
	}
	return nil
}
