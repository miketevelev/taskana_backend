package auth_postgres_repository

import (
	"context"
	"fmt"

	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (r *AuthRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO taskana.users (first_name, last_name, email, 
password_hash, timezone)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, version, first_name, last_name, email, password_hash, 
timezone, created_at, updated_at
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		user.Timezone,
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
		return domain.User{}, fmt.Errorf("scan user from db: %w", err)
	}

	userDomain := userDomainFromModel(userModel)

	return userDomain, nil
}

func (r *AuthRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, first_name, last_name, email, password_hash, 
timezone, created_at, updated_at
		FROM taskana.users
		WHERE email = $1
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		email,
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
		return domain.User{}, fmt.Errorf("scan user from db: %w", err)
	}

	userDomain := userDomainFromModel(userModel)

	return userDomain, nil
}
