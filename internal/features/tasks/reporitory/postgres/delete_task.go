package tasks_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
)

func (r *TaskRepository) DeleteTask(
	ctx context.Context,
	userID uuid.UUID,
	taskID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var projectID *uuid.UUID
	var position int

	err = tx.QueryRow(
		ctx, `
		SELECT project_id, position 
		FROM taskana.tasks 
		WHERE id = $1 AND user_id = $2`,
		taskID, userID,
	).Scan(&projectID, &position)

	if err != nil {
		return fmt.Errorf("failed to find task to delete: %w", err)
	}

	shiftQuery := `
		UPDATE taskana.tasks 
		SET position = position - 1
		WHERE user_id = $1 
		  AND project_id IS NOT DISTINCT FROM $2 
		  AND position > $3
	`
	_, err = tx.Exec(ctx, shiftQuery, userID, projectID, position)
	if err != nil {
		return fmt.Errorf("failed to shift positions: %w", err)
	}

	deleteQuery := `DELETE FROM taskana.tasks WHERE id = $1 AND user_id = $2`
	cmdTag, err := tx.Exec(ctx, deleteQuery, taskID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("task not found: %w", core_errors.ErrNotFound)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
