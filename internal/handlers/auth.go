package handlers

import (
	"fmt"
	"net/http"

	"github.com/Archiker-715/expense-tracker-api/internal/auth"
	"github.com/Archiker-715/expense-tracker-api/internal/entity"
	"github.com/Archiker-715/expense-tracker-api/internal/errs"
	"github.com/Archiker-715/expense-tracker-api/internal/repository/pg"
	"github.com/Archiker-715/expense-tracker-api/pkg/httpserver"
)

type AuthHandler struct {
	auth *auth.AuthService
}

func NewAuthHadler(repo *pg.AuthRepository) *AuthHandler {
	return &AuthHandler{auth: auth.NewAuthService(repo)}
}

func (a *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var user entity.UserAuthRegistration
	if err := httpserver.JsonDecode(w, r, &user, 0); err != nil {
		return
	}

	accessData, err := a.auth.Authorization(user)
	if err != nil {
		errs.WriteError(w, 0, http.StatusInternalServerError, fmt.Sprintf("Authorization error: %v", err))
	}

	if err := httpserver.JsonEncode(w, accessData, 0); err != nil {
		return
	}
}

func (a *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user entity.UserAuthRegistration
	if err := httpserver.JsonDecode(w, r, &user, 0); err != nil {
		return
	}

	err := a.auth.Registration(user)
	if err != nil {
		errs.WriteError(w, 0, http.StatusInternalServerError, fmt.Sprintf("Registration error: %v", err))
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Successfully sign Up! Now you can sign In via entered login & password")
}
