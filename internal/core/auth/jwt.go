package core_auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
)

type JWTConfig struct {
	Secret          string        `envconfig:"SECRET" required:"true"`
	AccessTokenTTL  time.Duration `envconfig:"ACCESS_TTL" default:"15m"`
	RefreshTokenTTL time.Duration `envconfig:"REFRESH_TTL" default:"720h"`
}

func NewJWTConfig() (JWTConfig, error) {
	var config JWTConfig
	if err := envconfig.Process("JWT", &config); err != nil {
		return JWTConfig{}, fmt.Errorf("load jwt config: %w", err)
	}
	return config, nil
}

func NewJWTConfigMust() JWTConfig {
	config, err := NewJWTConfig()
	if err != nil {
		panic(fmt.Errorf("jwt config: %w", err))
	}
	return config
}

type AccessClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	config JWTConfig
}

func NewTokenManager(config JWTConfig) *TokenManager {
	return &TokenManager{config: config}
}

func (m *TokenManager) GenerateAccessToken(userID uuid.UUID) (
	string,
	time.Time,
	error,
) {
	expiresAt := time.Now().UTC().Add(m.config.AccessTokenTTL)
	claims := AccessClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(m.config.Secret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign access token: %w", err)
	}

	return signed, expiresAt, nil
}

func (m *TokenManager) ParseAccessToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&AccessClaims{},
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf(
					"unexpected signing method: %w",
					core_errors.ErrUnauthorized,
				)
			}
			return []byte(m.config.Secret), nil
		},
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf(
			"parse access token: %w", core_errors.ErrUnauthorized,
		)
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return uuid.Nil, fmt.Errorf(
			"invalid access token: %w", core_errors.ErrUnauthorized,
		)
	}

	return claims.UserID, nil
}

func (m *TokenManager) RefreshTokenTTL() time.Duration {
	return m.config.RefreshTokenTTL
}

func (m *TokenManager) AccessTokenTTLSeconds() int64 {
	return int64(m.config.AccessTokenTTL.Seconds())
}

func GenerateRefreshToken(userID uuid.UUID) (string, error) {
	return userID.String() + "." + uuid.NewString() + uuid.NewString(), nil
}

func ParseRefreshTokenUserID(refreshToken string) (uuid.UUID, error) {
	parts := strings.SplitN(refreshToken, ".", 2)
	if len(parts) != 2 {
		return uuid.Nil, fmt.Errorf(
			"invalid refresh token format: %w", core_errors.ErrUnauthorized,
		)
	}
	return uuid.Parse(parts[0])
}
