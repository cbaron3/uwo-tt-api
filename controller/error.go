package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// HTTPError encapsulate http error fields
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

// NewError create a new http error and print message for debugging
func NewError(w http.ResponseWriter, status int, err error, logMsg string) http.ResponseWriter {
	fmt.Println(logMsg)

	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(er)

	// Return modified response
	return w
}
