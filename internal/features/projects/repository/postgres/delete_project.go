package projects_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
)

func (r *ProjectRepository) DeleteProject(
	ctx context.Context,
	userID uuid.UUID,
	projectID uuid.UUID,
) error {
	// todo: area must deleted all tasks inside
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		DELETE FROM taskana.projects
		WHERE id = $1 AND user_id = $2;
	`

	cmdTag, err := r.pool.Exec(ctx, query, projectID, userID)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf(
			"no project found with id '%s': %w",
			projectID,
			core_errors.ErrNotFound,
		)
	}

	return nil
}
