package checklists_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
	core_postgres_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool"
)

func (r *ChecklistsRepository) DeleteChecklist(
	ctx context.Context,
	userID uuid.UUID,
	checklistID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var taskID *uuid.UUID
	var position int

	deleteQuery := `
		DELETE FROM taskana.checklists 
		WHERE id = $1 AND user_id = $2
		RETURNING task_id, position
	`
	err = tx.QueryRow(ctx, deleteQuery, checklistID, userID).Scan(
		&taskID, &position,
	)
	if err != nil {
		if errors.Is(
			err, core_postgres_pool.ErrNoRows,
		) || err.Error() == "no rows in result set" {
			return fmt.Errorf(
				"checklist not found: %w", core_errors.ErrNotFound,
			)
		}
		return fmt.Errorf("failed to delete checklist: %w", err)
	}

	shiftQuery := `
		UPDATE taskana.checklists
		SET 
		    position = position - 1, 
		    version = version + 1, 
		    updated_at = NOW()
		WHERE user_id = $1
		  AND task_id IS NOT DISTINCT FROM $2
		  AND position > $3
	`
	_, err = tx.Exec(ctx, shiftQuery, userID, taskID, position)
	if err != nil {
		return fmt.Errorf("failed to shift positions: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
