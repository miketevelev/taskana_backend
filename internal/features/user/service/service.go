package user_service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type UsersService struct {
	userRepository UserRepository
	tokenManager   *core_auth.TokenManager
}

type UserRepository interface {
	GetUser(
		ctx context.Context,
		userID uuid.UUID,
	) (domain.User, error)

	GetUserByID(
		ctx context.Context,
		id uuid.UUID,
	) (domain.User, error)

	PatchUser(
		ctx context.Context,
		userID uuid.UUID,
		user domain.User,
	) (domain.User, error)

	ChangePassword(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	SaveRefreshToken(
		ctx context.Context,
		userID uuid.UUID,
		tokenHash string,
		expiresAt time.Time,
		userAgent *string,
	) error

	CheckEmail(
		ctx context.Context,
		email string,
	) error
}

func NewUsersService(
	userRepository UserRepository,
	tokenManager *core_auth.TokenManager,
) *UsersService {
	return &UsersService{
		userRepository: userRepository,
		tokenManager:   tokenManager,
	}
}

func (s *UsersService) issueTokenPair(
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
	if err := s.userRepository.SaveRefreshToken(
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
