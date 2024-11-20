package utils

import (
	"net/http"
	"time"
)

type Cookie struct {
	Name     string
	Value    string
	Expires  time.Time
	HttpOnly bool
	Secure   bool
	Path     string
	SameSite http.SameSite
}

func NewCookie(name, value string, expires time.Time, httpOnly, secure bool) *Cookie {
	return &Cookie{
		Name:     name,
		Value:    value,
		Expires:  expires,
		HttpOnly: httpOnly,
		Secure:   secure,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
}

func (c *Cookie) SetCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     c.Name,
		Value:    c.Value,
		Expires:  c.Expires,
		HttpOnly: c.HttpOnly,
		Secure:   c.Secure,
		Path:     c.Path,
		SameSite: c.SameSite,
	})
}

func (c *Cookie) GetCookie(r *http.Request) (*http.Cookie, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, err
	}
	return cookie, nil
}
