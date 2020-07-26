package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// HTTPError example
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

// NewError example
func NewError(w http.ResponseWriter, status int, err error, logMsg string) http.ResponseWriter {
	fmt.Println(logMsg)

	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(er)

	return w
}
