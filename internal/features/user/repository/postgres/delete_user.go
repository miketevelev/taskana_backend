package user_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
)

func (r *UserRepository) DeleteUser(
	ctx context.Context,
	userID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tag, err := r.pool.Exec(
		ctx,
		`DELETE FROM taskana.users WHERE id = $1`,
		userID,
	)
	if err != nil {
		return fmt.Errorf("delete user repository: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return core_errors.ErrNotFound
	}

	return nil
}
