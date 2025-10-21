package main

import (
	"encoding/json"
	"net/http"
)

func toApiResponse(w http.ResponseWriter, statusCode int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := apiResponse{
		message,
		data,
	}

	json.NewEncoder(w).Encode(resp)
}