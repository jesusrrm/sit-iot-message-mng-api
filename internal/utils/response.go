package utils

import (
	"net/http"
	"encoding/json"
)

// Response represents the standard API response structure
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SendResponse sends a JSON response with the given status code
func SendResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := Response{
		Status:  http.StatusText(statusCode),
		Message: message,
		Data:    data,
	}
	json.NewEncoder(w).Encode(response)
}

// SendError sends a JSON error response
func SendError(w http.ResponseWriter, statusCode int, message string) {
	SendResponse(w, statusCode, message, nil)
}