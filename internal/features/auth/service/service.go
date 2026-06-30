package auth_service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type AuthService struct {
	authRepository AuthRepository
	tokenManager   *core_auth.TokenManager
}

type AuthRepository interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	GetUserByEmail(
		ctx context.Context,
		email string,
	) (domain.User, error)

	GetRefreshToken(
		ctx context.Context,
		tokenHash string,
	) (uuid.UUID, time.Time, error)

	DeleteRefreshToken(
		ctx context.Context,
		tokenHash string,
	) error

	SaveRefreshToken(
		ctx context.Context,
		userID uuid.UUID,
		tokenHash string,
		expiresAt time.Time,
		userAgent *string,
	) error

	DeleteAllRefreshTokens(
		ctx context.Context,
		userID uuid.UUID,
	) error
}

func NewAuthService(
	authRepository AuthRepository,
	tokenManager *core_auth.TokenManager,
) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		tokenManager:   tokenManager,
	}
}

func (s *AuthService) issueTokenPair(
	ctx context.Context,
	userID uuid.UUID,
	userAgent *string,
) (domain.TokenPair, error) {
	accessToken, _, err := s.tokenManager.GenerateAccessToken(userID)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf(
			"generating access token failed: %w", err,
		)
	}

	refreshToken, err := core_auth.GenerateRefreshToken(userID)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf(
			"generating refresh token failed: %w", err,
		)
	}

	tokenHash, err := core_auth.HashToken(refreshToken)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf(
			"hashing refresh token failed: %w", err,
		)
	}

	expiresAt := time.Now().UTC().Add(s.tokenManager.RefreshTokenTTL())
	if err := s.authRepository.SaveRefreshToken(
		ctx,
		userID,
		tokenHash,
		expiresAt,
		userAgent,
	); err != nil {
		return domain.TokenPair{}, fmt.Errorf(
			"saving refresh token failed: %w", err,
		)
	}

	return domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.tokenManager.AccessTokenTTLSeconds(),
	}, nil
}
