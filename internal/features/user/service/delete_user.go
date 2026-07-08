package user_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *UsersService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	if err := s.userRepository.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("delete user service: %w", err)
	}

	return nil
}
