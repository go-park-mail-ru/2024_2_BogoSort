package utils

import (
	"encoding/json"
	"net/http"
)

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
		if _, err := w.Write(errJSON); err != nil {
			SendErrorResponse(w, http.StatusInternalServerError, "Failed to encode response")
		}
	}
}

func SendJSONResponse(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to encode response")
	}
}
