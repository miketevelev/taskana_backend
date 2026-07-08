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
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var areaID *uuid.UUID
	var position int

	err = tx.QueryRow(
		ctx, `
		SELECT area_id, position 
		FROM taskana.projects 
		WHERE id = $1 AND user_id = $2`,
		projectID, userID,
	).Scan(&areaID, &position)

	if err != nil {
		return fmt.Errorf("failed to find project to delete: %w", err)
	}

	shiftQuery := `
		UPDATE taskana.projects 
		SET position = position - 1
		WHERE user_id = $1 
		  AND area_id IS NOT DISTINCT FROM $2 
		  AND position > $3
	`
	_, err = tx.Exec(ctx, shiftQuery, userID, areaID, position)
	if err != nil {
		return fmt.Errorf("failed to shift positions: %w", err)
	}

	deleteQuery := `DELETE FROM taskana.projects WHERE id = $1 AND user_id = $2`
	cmdTag, err := tx.Exec(ctx, deleteQuery, projectID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("project not found: %w", core_errors.ErrNotFound)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
