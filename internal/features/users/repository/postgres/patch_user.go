package user_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
	core_postgres_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool"
)

func (r *UserRepository) CheckEmail(
	ctx context.Context,
	email string,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM taskana.users 
			WHERE email = $1
		);
	`

	var exists bool
	err := r.pool.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check email existence: %w", err)
	}

	if exists {
		return core_errors.ErrAlreadyExists
	}

	return nil

	return nil
}

func (r *UserRepository) PatchUser(
	ctx context.Context,
	userID uuid.UUID,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE taskana.users
	SET
		first_name = $1,
		last_name = $2,
		email = $3,
		timezone = $4,
		updated_at = NOW(),
		version=version + 1
	WHERE id = $5 AND version = $6
	RETURNING id, version, first_name, last_name, email, password_hash, 
timezone, created_at, updated_at
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Timezone,
		userID,
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
			return domain.User{},
				fmt.Errorf(
					"users with id='%s' concurrently accessed or not found: %w",
					userID,
					core_errors.ErrConflict,
				)
		}
		return domain.User{}, fmt.Errorf("patch users repository: %w", err)
	}

	userDomain := userDomainFromModel(userModel)

	return userDomain, nil
}
