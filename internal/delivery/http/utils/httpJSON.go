package utils

import (
	"io"
	"net/http"

	"github.com/mailru/easyjson"
)

func writeJSON(w http.ResponseWriter, v easyjson.Marshaler, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	writer := io.Writer(w)
	_, err := easyjson.MarshalToWriter(v, writer)
	return err
}

func readJSON(r *http.Request, v easyjson.Unmarshaler) error {
	if err := easyjson.UnmarshalFromReader(r.Body, v); err != nil {
		return err
	}
	return r.Body.Close()
}
