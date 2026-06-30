package auth_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
	core_postgres_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool"
)

func (r *AuthRepository) SaveRefreshToken(
	ctx context.Context,
	userID uuid.UUID,
	tokenHash string,
	expiresAt time.Time,
	userAgent *string,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO taskana.refresh_tokens (user_id, token_hash, expires_at, user_agent)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.pool.Exec(ctx, query, userID, tokenHash, expiresAt, userAgent)
	if err != nil {
		return fmt.Errorf("save refresh token: %w", err)
	}

	return nil
}

func (r *AuthRepository) GetRefreshToken(
	ctx context.Context,
	tokenHash string,
) (uuid.UUID, time.Time, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT user_id, expires_at
		FROM taskana.refresh_tokens
		WHERE token_hash = $1
	`

	var userID uuid.UUID
	var expiresAt time.Time
	err := r.pool.QueryRow(ctx, query, tokenHash).Scan(&userID, &expiresAt)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return uuid.Nil, time.Time{}, core_errors.ErrNotFound
		}
		return uuid.Nil, time.Time{}, fmt.Errorf("get refresh token: %w", err)
	}

	return userID, expiresAt, nil
}

func (r *AuthRepository) DeleteRefreshToken(
	ctx context.Context,
	tokenHash string,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `DELETE FROM taskana.refresh_tokens WHERE token_hash = $1`
	tag, err := r.pool.Exec(ctx, query, tokenHash)
	if err != nil {
		return fmt.Errorf("delete refresh token: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return core_errors.ErrNotFound
	}
	return nil
}

func (r *AuthRepository) DeleteAllRefreshTokens(
	ctx context.Context,
	userID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `DELETE FROM taskana.refresh_tokens WHERE user_id = $1`
	_, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("delete all refresh tokens: %w", err)
	}
	return nil
}
