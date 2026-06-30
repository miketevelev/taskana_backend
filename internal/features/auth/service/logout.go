package auth_service

import (
	"context"
	"fmt"

	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
)

func (s *AuthService) Logout(
	ctx context.Context,
	refreshToken string,
) error {
	tokenHash, err := core_auth.HashToken(refreshToken)
	if err != nil {
		return fmt.Errorf("hash refresh token: %w", err)
	}

	if err := s.authRepository.DeleteRefreshToken(ctx, tokenHash); err != nil {
		return fmt.Errorf("delete refresh token: %w", err)
	}

	return nil
}
