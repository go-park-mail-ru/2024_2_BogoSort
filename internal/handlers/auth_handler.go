package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"emporium/internal/storage"

	"github.com/golang-jwt/jwt/v5"

	"log"
	"errors"
)

type LoginCredentials struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

const (
	sessionTime = 24 * time.Hour
)

var jwtKey = []byte("your_secret_key")

type AuthHandler struct {
	UserStorage    *storage.UserStorage
	SessionStorage *storage.SessionStorage
}

func (ah *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var credentials LoginCredentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := ah.UserStorage.CreateUser(credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   tokenString,
		Expires: time.Now().Add(sessionTime),
	}

	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func (ah *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var credentials LoginCredentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := ah.UserStorage.ValidateUserByEmailAndPassword(credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   tokenString,
		Expires: time.Now().Add(sessionTime),
	}

	http.SetCookie(w, cookie)

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func (ah *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "No active session", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Failed to retrieve cookie", http.StatusBadRequest)
		return
	}

	// Invalidate the session
	cookie.Expires = time.Now().Add(-1 * time.Hour)
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("Unexpected signing method: %v", token.Header["alg"])
				return nil, errors.New("unexpected signing method")
			}
			return jwtKey, nil
		})

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			r.Header.Set("User", claims["email"].(string))
			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

// func LoginHandler(w http.ResponseWriter, r *http.Request) {
// 	var credentials LoginCredentials
// 	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	userStorage := storage.NewUserStorage()
// 	user, err := userStorage.ValidateUserByEmailAndPassword(credentials.Email, credentials.Password)
// 	if err != nil {
// 		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
// 		return
// 	}

// 	expirationTime := time.Now().Add(sessionTime)
// 	claims := &jwt.StandardClaims{
// 		Subject:   user.Email,
// 		ExpiresAt: expirationTime.Unix(),
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString(jwtKey)
// 	if err != nil {
// 		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
// 		return
// 	}

// 	sessionStorage := storage.NewSessionStorage()
// 	sessionStorage.AddSession(tokenString)

// 	http.SetCookie(w, &http.Cookie{
// 		Name:    "token",
// 		Value:   tokenString,
// 		Expires: expirationTime,
// 	})

// 	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
// }

// func LogoutHandler(w http.ResponseWriter, r *http.Request) {
// 	c, err := r.Cookie("token")
// 	if err != nil {
// 		if err == http.ErrNoCookie {
// 			http.Error(w, "No active session", http.StatusUnauthorized)
// 			return
// 		}
// 		http.Error(w, "Failed to retrieve cookie", http.StatusBadRequest)
// 		return
// 	}
// 	tokenString := c.Value

// 	sessionStorage := storage.NewSessionStorage()
// 	sessionStorage.RemoveSession(tokenString)

// 	http.SetCookie(w, &http.Cookie{
// 		Name:    "token",
// 		Value:   "",
// 		Expires: time.Now(),
// 	})

// 	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
// }
