package main

import (
	"encoding/json"
	"net/http"
)

type apiResponse struct {
	Data any `json:"data"`
}

type apiError struct {
	Error string `json:"error"`
}

func (app application) jsonResponse(w http.ResponseWriter, status int, data any) error {
	return writeJSON(w, status, apiResponse{Data: data})
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, dest any) error {
	maxBytes := 1_048_576 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(dest)
}

func writeJSONError(w http.ResponseWriter, status int, message string) error {
	return writeJSON(w, status, apiError{
		Error: message,
	})
}
