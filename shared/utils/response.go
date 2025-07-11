package utils

import (
	"dailytrackr/shared/dto"
	"encoding/json"
	"net/http"
)

// SendSuccessResponse sends a successful JSON response
func SendSuccessResponse(w http.ResponseWriter, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := dto.Response{
		Success: true,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

// SendErrorResponse sends an error JSON response
func SendErrorResponse(w http.ResponseWriter, statusCode int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := dto.ErrorResponse{
		Success: false,
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	json.NewEncoder(w).Encode(response)
}

// SendCreatedResponse sends a 201 Created response
func SendCreatedResponse(w http.ResponseWriter, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := dto.Response{
		Success: true,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

// SendBadRequestResponse sends a 400 Bad Request response
func SendBadRequestResponse(w http.ResponseWriter, message string, err error) {
	SendErrorResponse(w, http.StatusBadRequest, message, err)
}

// SendUnauthorizedResponse sends a 401 Unauthorized response
func SendUnauthorizedResponse(w http.ResponseWriter, message string) {
	SendErrorResponse(w, http.StatusUnauthorized, message, nil)
}

// SendNotFoundResponse sends a 404 Not Found response
func SendNotFoundResponse(w http.ResponseWriter, message string) {
	SendErrorResponse(w, http.StatusNotFound, message, nil)
}

// SendInternalServerErrorResponse sends a 500 Internal Server Error response
func SendInternalServerErrorResponse(w http.ResponseWriter, message string, err error) {
	SendErrorResponse(w, http.StatusInternalServerError, message, err)
}
