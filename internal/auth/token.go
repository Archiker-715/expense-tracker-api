package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/Archiker-715/expense-tracker-api/internal/entity"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (a *AuthService) generateToken(userID uuid.UUID) (entity.AuthResponse, error) {
	expTokenTime := time.Now().Add(time.Hour * 1).Unix()

	claims := jwt.MapClaims{
		"user_id": userID,
		"iat":     time.Now().Unix(),
		"exp":     expTokenTime,
		"jti":     uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	outToken, err := token.SignedString(secretKey)
	if err != nil {
		return entity.AuthResponse{}, errors.New("signing token error")
	}

	out := entity.AuthResponse{
		Token:     outToken,
		ExpiresIn: int(expTokenTime - time.Now().Unix()),
	}

	return out, nil
}

func ParseToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
}
