package user_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
	core_postgres_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool"
)

func (r *UserRepository) SaveRefreshToken(
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

func (r *UserRepository) ChangePassword(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE taskana.users
		SET
			password_hash = $1,
			updated_at = NOW(),
			version = version + 1
		WHERE id = $2 AND version = $3
		RETURNING id, version, first_name, last_name, email, password_hash, 
timezone, created_at, updated_at
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.PasswordHash,
		user.ID,
		user.Version,
	)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FirstName,
		&userModel.LastName,
		&userModel.Email,
		&userModel.PasswordHash,
		&userModel.Timezone,
		&userModel.CreatedAt,
		&userModel.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"users with id='%s' concurrently accessed or not found: %w",
				user.ID,
				core_errors.ErrConflict,
			)
		}
		return domain.User{}, fmt.Errorf("change password repository: %w", err)
	}

	userDomain := userDomainFromModel(userModel)

	return userDomain, nil
}

func (r *UserRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, first_name, last_name, email, password_hash, 
timezone, created_at, updated_at
		FROM taskana.users
		WHERE id = $1
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		id,
	)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FirstName,
		&userModel.LastName,
		&userModel.Email,
		&userModel.PasswordHash,
		&userModel.Timezone,
		&userModel.CreatedAt,
		&userModel.UpdatedAt,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("scan users from db: %w", err)
	}

	userDomain := userDomainFromModel(userModel)

	return userDomain, nil
}
