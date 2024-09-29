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
	json.NewEncoder(w).Encode(AuthErrResponse{Code: code, Status: status})
}

func SendJSONResponse(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
