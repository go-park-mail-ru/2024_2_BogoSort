package http

import (
    "net/http"
    "time"

    "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/http/utils"
    "github.com/google/uuid"
    "github.com/gorilla/mux"
)

type CSRFEndpoint struct {
	csrfTokenUtil  *utils.CryptToken
	sessionManager *utils.SessionManager
}

func NewCSRFEndpoint(csrfTokenUtil *utils.CryptToken, sessionManager *utils.SessionManager) *CSRFEndpoint {
    return &CSRFEndpoint{
        csrfTokenUtil:  csrfTokenUtil,
        sessionManager: sessionManager,
    }
}

func (c *CSRFEndpoint) Configure(router *mux.Router) {
	router.HandleFunc("/api/v1/csrf-token", c.Get).Methods(http.MethodGet)
}

// Get handles the HTTP request to retrieve a CSRF token.
// @Summary Retrieve CSRF Token
// @Description This endpoint checks for a session ID in the request cookies and retrieves the user ID from the session manager. If both are valid, it generates a CSRF token using the session ID and user ID, and sends it back in the response header. If any step fails, it responds with an appropriate error message.
// @Tags CSRF
// @Accept json
// @Produce json
// @Success 200 {string} string "CSRF Token"
// @Failure 401 {object} utils.ErrResponse "Unauthorized"
// @Failure 500 {object} utils.ErrResponse "Failed to create CSRF token"
// @Router /api/v1/csrf-token [get]
func (c *CSRFEndpoint) Get(w http.ResponseWriter, r *http.Request) {
	sessionID, err := r.Cookie("session_id")
	if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    userID, err := c.sessionManager.GetUserID(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    token, err := c.csrfTokenUtil.Create(uuid.MustParse(sessionID.Value), userID, time.Now().Add(24*time.Hour).Unix())
    if err != nil {
        http.Error(w, "Failed to create CSRF token", http.StatusInternalServerError)
        return
    }

    w.Header().Set("X-CSRF-Token", token)
    w.WriteHeader(http.StatusOK)
}