package auth_service

import (
	"context"
	"fmt"
	"time"

	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
)

func (s *AuthService) Refresh(
	ctx context.Context,
	refreshToken string,
	userAgent *string,
) (domain.TokenPair, error) {
	tokenHash, err := core_auth.HashToken(refreshToken)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf(
			"failed to hash refresh token: %w", err,
		)
	}

	userID, expiresAt, err := s.authRepository.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		if replayUserID, replayErr := core_auth.ParseRefreshTokenUserID(refreshToken); replayErr == nil {
			_ = s.authRepository.DeleteAllRefreshTokens(ctx, replayUserID)
		}
		return domain.TokenPair{}, fmt.Errorf(
			"invalid refresh token: %w", core_errors.ErrUnauthorized,
		)
	}

	if time.Now().UTC().After(expiresAt) {
		_ = s.authRepository.DeleteRefreshToken(ctx, tokenHash)
		return domain.TokenPair{}, fmt.Errorf(
			"refresh token expired: %w", core_errors.ErrUnauthorized,
		)
	}

	if err := s.authRepository.DeleteRefreshToken(ctx, tokenHash); err != nil {
		return domain.TokenPair{}, fmt.Errorf("rotate refresh token: %w", err)
	}

	return s.issueTokenPair(ctx, userID, userAgent)
}
