package utils

import (
	"errors"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/config"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken            = errors.New("invalid token")
)

func InitJWT() {
	jwtKey = []byte(config.GetJWTSecretKey())
	if len(jwtKey) == 0 {
		panic("JWT secret key is not set")
	}
}

func CreateToken(email string) (string, error) {
	expirationTime := time.Now().Add(config.GetJWTExpirationTime())
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    config.GetJWTIssuer(),
		Subject:   email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}

func CreateRefreshToken(email string) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    config.GetJWTIssuer(),
		Subject:   email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}

func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}

		return jwtKey, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims.Subject, nil
	}

	return "", ErrInvalidToken
}
