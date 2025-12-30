package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type tokenContextKey struct{}
type UserContextKey struct{}

var TokenCtxKey = tokenContextKey{}
var UserCtxKey = UserContextKey{}

func TokenToContext(ctx context.Context, token *jwt.Token) context.Context {
	return context.WithValue(ctx, TokenCtxKey, token)
}

func TokenFromContext(ctx context.Context) (*jwt.Token, bool) {
	token, ok := ctx.Value(TokenCtxKey).(*jwt.Token)
	return token, ok
}

func UserToContext(ctx context.Context, userId uuid.UUID) context.Context {
	return context.WithValue(ctx, UserCtxKey, userId)
}

func UserFromContext(ctx context.Context) (uuid.UUID, bool) {
	user, ok := ctx.Value(UserCtxKey).(uuid.UUID)
	return user, ok
}
