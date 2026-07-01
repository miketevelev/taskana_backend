package auth_service

import (
	"context"
	"fmt"

	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
)

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
	userAgent *string,
) (domain.TokenPair, domain.User, error) {
	user, err := s.authRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.TokenPair{}, domain.User{}, fmt.Errorf(
			"get user by email: %v: %w", core_errors.ErrUnauthorized, err,
		)
	}

	ok, err := core_auth.VerifyPassword(password, user.PasswordHash)
	if err != nil || !ok {
		return domain.TokenPair{}, domain.User{},
			fmt.Errorf(
				"invalid password: %v: %w", core_errors.ErrUnauthorized, err,
			)
	}

	tokens, err := s.issueTokenPair(ctx, user.ID, userAgent)
	if err != nil {
		return domain.TokenPair{}, domain.User{}, fmt.Errorf(
			"issue token pair: %v: %w", core_errors.ErrUnauthorized, err,
		)
	}

	return tokens, user, nil
}
