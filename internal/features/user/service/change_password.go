package user_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
)

func (s *UsersService) ChangePassword(
	ctx context.Context,
	userID uuid.UUID,
	oldPassword string,
	newPassword string,
	userAgent *string,
) (domain.TokenPair, domain.User, error) {
	if oldPassword == newPassword {
		return domain.TokenPair{}, domain.User{}, fmt.Errorf(
			"new password cannot be the same as the old password: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if len(newPassword) < 8 {
		return domain.TokenPair{}, domain.User{}, fmt.Errorf(
			"password must be at least 8 characters long: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	user, err := s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return domain.TokenPair{}, domain.User{}, fmt.Errorf(
			"get user by ID: %w", err,
		)
	}

	ok, err := core_auth.VerifyPassword(oldPassword, user.PasswordHash)
	if err != nil {
		return domain.TokenPair{}, domain.User{}, fmt.Errorf(
			"verify password failed: %s: %w", err.Error(),
			core_errors.ErrUnauthorized,
		)
	}
	if !ok {
		return domain.TokenPair{}, domain.User{}, fmt.Errorf(
			"incorrect old password: %w", core_errors.ErrUnauthorized,
		)
	}

	hash, err := core_auth.HashPassword(newPassword)
	if err != nil {
		return domain.TokenPair{},
			domain.User{},
			fmt.Errorf("hashing password failed: %w", err)
	}
	user.PasswordHash = hash

	user, err = s.userRepository.ChangePassword(ctx, user)
	if err != nil {
		return domain.TokenPair{}, domain.User{}, fmt.Errorf(
			"change password: %w", err,
		)
	}

	tokens, err := s.issueTokenPair(ctx, user.ID, userAgent)
	if err != nil {
		return domain.TokenPair{}, domain.User{}, fmt.Errorf(
			"issuing token pair: %w", err,
		)
	}

	return tokens, user, nil
}
