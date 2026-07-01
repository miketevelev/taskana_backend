package user_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
)

func (s *UsersService) PatchUser(
	ctx context.Context,
	userID uuid.UUID,
	patch domain.UserPatch,
) (domain.User, error) {
	user, err := s.userRepository.GetUser(ctx, userID)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user service: %w", err)
	}

	oldEmail := user.Email

	if patch.Email.Value != nil && *patch.Email.Value == oldEmail {
		return domain.User{}, fmt.Errorf(
			"email address is the same as old email: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if err := user.ApplyPatch(patch); err != nil {
		return domain.User{}, fmt.Errorf("patch user service: %w", err)
	}

	if user.Email != oldEmail {
		if err := s.userRepository.CheckEmail(ctx, user.Email); err != nil {
			if errors.Is(err, core_errors.ErrAlreadyExists) {
				return domain.User{}, fmt.Errorf(
					"email already exists: %w", err,
				)
			}
			return domain.User{}, fmt.Errorf("failed to check email: %w", err)
		}
	}

	patchedUser, err := s.userRepository.PatchUser(ctx, userID, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("patch user service: %w", err)
	}

	return patchedUser, nil
}
