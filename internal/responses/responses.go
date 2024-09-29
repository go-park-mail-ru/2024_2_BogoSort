package responses

import (
	"encoding/json"
	"net/http"
)

type AuthResponse struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

type AuthErrResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

func SendErrorResponse(w http.ResponseWriter, code int, status string) {
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(AuthErrResponse{Code: code, Status: status})
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func SendJSONResponse(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
