package user_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type UsersService struct {
	userRepository UserRepository
}

type UserRepository interface {
	GetUser(
		ctx context.Context,
		userID uuid.UUID,
	) (domain.User, error)

	PatchUser(
		ctx context.Context,
		userID uuid.UUID,
		user domain.User,
	) (domain.User, error)

	CheckEmail(
		ctx context.Context,
		email string,
	) error
}

func NewUsersService(userRepository UserRepository) *UsersService {
	return &UsersService{
		userRepository: userRepository,
	}
}
