package areas_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
)

func (r *AreasRepository) DeleteArea(
	ctx context.Context,
	userID uuid.UUID,
	areaID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		DELETE FROM taskana.areas 
		WHERE id = $1 AND user_id = $2;
	`

	cmdTag, err := r.pool.Exec(ctx, query, areaID, userID)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf(
			"no task found with id '%d': %w",
			areaID,
			core_errors.ErrNotFound,
		)
	}

	return nil
}
