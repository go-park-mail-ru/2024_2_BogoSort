package responses

import (
	"encoding/json"
	"net/http"
)

type AuthResponse struct {
	Email string `json:"email"`
	SessionID string `json:"session_id"`
	IsAuth    bool   `json:"is_auth"`
}

type ErrResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

func SendErrorResponse(w http.ResponseWriter, code int, status string) {
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(ErrResponse{Code: code, Status: status})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errResponse := ErrResponse{
			Code:   http.StatusInternalServerError,
			Status: "Failed to encode response",
		}
		errJSON, _ := json.Marshal(errResponse)
		w.Write(errJSON)
	}
}

func SendJSONResponse(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(payload)

	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to encode response")
	}
}
