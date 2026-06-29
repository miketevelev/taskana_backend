package core_auth

import (
	"context"

	"github.com/google/uuid"
)

type contextKey struct{}

var userIDKey = contextKey{}

func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	val, ok := ctx.Value(userIDKey).(uuid.UUID)
	return val, ok && val != uuid.Nil
}

func MustUserIDFromContext(ctx context.Context) uuid.UUID {
	userID, ok := UserIDFromContext(ctx)
	if !ok {
		panic("user_id not found in context")
	}
	return userID
}
