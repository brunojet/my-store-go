package base

import (
	"encoding/json"
	"net/http"
)

type BaseHandler struct{}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h BaseHandler) WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (h BaseHandler) WriteError(w http.ResponseWriter, status int, message string) {
	h.WriteJSON(w, status, ErrorResponse{Error: message})
}
