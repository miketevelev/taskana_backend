package auth_service

import (
	"context"
	"fmt"

	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *AuthService) Register(
	ctx context.Context,
	firstName string,
	lastName string,
	email string,
	password string,
	timezone string,
) (domain.TokenPair, domain.User, error) {
	user := domain.NewUserUninitialized(
		firstName,
		lastName,
		email,
		password,
		timezone,
	)
	if err := user.Validate(); err != nil {
		return domain.TokenPair{},
			domain.User{},
			fmt.Errorf("users validation failed: %w", err)
	}

	hash, err := core_auth.HashPassword(password)
	if err != nil {
		return domain.TokenPair{},
			domain.User{},
			fmt.Errorf("hashing password failed: %w", err)
	}
	user.PasswordHash = hash

	user, err = s.authRepository.CreateUser(ctx, user)
	if err != nil {
		return domain.TokenPair{}, domain.User{}, fmt.Errorf(
			"users creation failed: %w", err,
		)
	}

	tokens, err := s.issueTokenPair(ctx, user.ID, nil)
	if err != nil {
		return domain.TokenPair{}, domain.User{}, fmt.Errorf(
			"issuing token pair: %w", err,
		)
	}

	return tokens, user, nil
}
