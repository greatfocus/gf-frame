package server

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/greatfocus/gf-sframe/crypt"
)

// Response data
type Response struct {
	Payload string `json:"data,omitempty"`
}

// Success returns response as json
func Success(w http.ResponseWriter, statusCode int, data interface{}) {
	res := Response{Payload: crypt.Encrypt(data.(string), os.Args[2])}
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(res)
}

// Error returns error as json
func Error(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		Success(w, statusCode, struct {
			Error string `json:"error"`
		}{Error: err.Error()})
		return
	}
	Success(w, http.StatusBadRequest, nil)
}
