package responser

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func WriteJSON(w http.ResponseWriter, status int, response []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(response)))
	WriteStatus(w, status)
	_, err := w.Write(response)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func WriteError(w http.ResponseWriter, status int, err error) {
	errJSON, _ := json.Marshal(ErrorResponse{err.Error()})
	WriteStatus(w, status)
	_, _ = w.Write(errJSON)
}
