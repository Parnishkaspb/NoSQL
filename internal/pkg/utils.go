package pkg

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func WriteApiResponse(w http.ResponseWriter, result any, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := ApiResponse{
		Message: message,
		Data:    result,
	}

	json.NewEncoder(w).Encode(response)
}
