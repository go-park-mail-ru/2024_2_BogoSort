package entity

import (
	"bytes"
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/random"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
	"regexp"
	_ "time"
	"unicode/utf8"
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
	Status       uint      `db:"status"`
}

func ValidatePassword(password string) error {
	switch {
	case utf8.RuneCountInString(password) < 8:
		return errors.New("пароль должен содержать не менее 8 символов")
	case utf8.RuneCountInString(password) > 32:
		return errors.New("пароль должен содержать не более 32 символов")
	case !regexp.MustCompile(`^[!@#$%^&*()_+\-=.,\w]+$`).MatchString(password):
		return errors.New("пароль может состоять из латинских букв, цифр и " +
			"специальных символов !@#$%^&*()_+\\-=")
	default:
		return nil
	}
}

func ValidateEmail(email string) error {
	re := regexp.MustCompile("^([a-z0-9!#$%&'*+\\\\/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+\\\\/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)$") // nolint: lll

	if !re.MatchString(email) {
		return errors.New("невалидная почта")
	}
	if len(email) > 256 {
		return errors.New("почта не может быть длиннее 256 символов")
	}
	return nil
}

func HashPassword(password string) (salt []byte, hash []byte, err error) {
	salt, err = random.Bytes(8)
	if err != nil {
		return nil, nil, errors.New("произошла непредвиденная ошибка")
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
		return errors.New("имя не может быть длиннее 30 символов")
	}
	return nil
}
