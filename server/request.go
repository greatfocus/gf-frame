package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/greatfocus/gf-frame/crypt"
)

// Request data
type Request struct {
	Payload string `json:"data,omitempty"`
}

// Success returns response as json
func GetPayload(w http.ResponseWriter, r *http.Request) (bool, []byte) {
	// check if the body is valid
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error(w, http.StatusBadRequest, err)
		return false, nil
	}

	// Get the Payload string and convert to requst struct
	req := Request{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		Error(w, http.StatusBadRequest, err)
		return false, nil
	}

	// decrypt the string and return byte
	payload := crypt.Decrypt(req.Payload, os.Args[2])
	return true, []byte(payload)
}
