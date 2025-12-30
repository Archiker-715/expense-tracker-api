package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Archiker-715/expense-tracker-api/internal/errs"
)

func JsonDecode(w http.ResponseWriter, r *http.Request, dest interface{}, internalCode int) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		errs.WriteError(w, internalCode, http.StatusBadRequest, fmt.Sprintf("invalid json: %v", err))
		return err
	}
	return nil
}

func JsonEncode(w http.ResponseWriter, dest interface{}, internalCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dest); err != nil {
		errs.WriteError(w, 0, http.StatusInternalServerError, fmt.Sprintf("build json err: %v", err))
		return err
	}
	return nil
}
