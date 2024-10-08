package responses

import (
	"encoding/json"
	"net/http"
)

type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	Email        string `json:"email"`
}

type ErrResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

func SendErrorResponse(w http.ResponseWriter, code int, status string) {
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(ErrResponse{Code: code, Status: status})

	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func SendJSONResponse(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(payload)

	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to encode response")
	}
}
