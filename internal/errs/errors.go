package errs

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
}

func WriteError(w http.ResponseWriter, internalCode, httpCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	errResp := errorResponse{
		Code:    internalCode,
		Message: message,
	}
	json.NewEncoder(w).Encode(errResp)
}
