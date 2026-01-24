package middleware

import (
	"net/http"
	"strings"

	"github.com/Archiker-715/expense-tracker-api/internal/auth"
	"github.com/Archiker-715/expense-tracker-api/internal/errs"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errs.WriteError(w, 0, http.StatusUnauthorized, "empty token")
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := auth.ParseToken(tokenStr)
		if err != nil || !token.Valid {
			errs.WriteError(w, 0, http.StatusUnauthorized, "invalid token")
			return
		}
		ctx := auth.TokenToContext(r.Context(), token)

		userId := getUUIDfromToken(token)
		if userId == uuid.Nil {
			errs.WriteError(w, 0, http.StatusUnauthorized, "couldn't extract uuid from token")
			return
		}
		ctx = auth.UserToContext(r.Context(), userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUUIDfromToken(token *jwt.Token) uuid.UUID {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil
	}

	subClaim, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil
	}

	userId, err := uuid.Parse(subClaim)
	if err != nil {
		return uuid.Nil
	}

	return userId
}
