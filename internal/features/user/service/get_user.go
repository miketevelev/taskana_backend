package user_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *UsersService) GetUser(
	ctx context.Context,
	userID uuid.UUID,
) (domain.User, error) {
	user, err := s.userRepository.GetUser(ctx, userID)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user from repository: %w", err)
	}

	return user, nil
}
