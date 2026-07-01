package user_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (r *UserRepository) GetUser(
	ctx context.Context,
	userID uuid.UUID,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, first_name, last_name, email, password_hash, 
timezone, created_at, updated_at
		FROM taskana.user
		WHERE id = $1;
		`

	row := r.pool.QueryRow(ctx, query, userID)

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
