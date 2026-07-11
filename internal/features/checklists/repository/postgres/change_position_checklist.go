package checklists_postgres_repository

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

func (r *ChecklistsRepository) ChangePosition(
	ctx context.Context,
	userID uuid.UUID,
	checklist domain.Checklist,
	oldPosition int,
) (domain.Checklist, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Checklist{}, err
	}
	defer tx.Rollback(ctx)

	var count int
	countQuery := `
		SELECT COUNT(*)
		FROM taskana.checklists
		WHERE user_id = $1 AND task_id IS NOT DISTINCT FROM $2
	`
	if err := tx.QueryRow(
		ctx, countQuery, userID, checklist.TaskID,
	).Scan(&count); err != nil {
		return domain.Checklist{}, fmt.Errorf(
			"failed to get checklist count: %w", err,
		)
	}

	newPos := checklist.Position
	if newPos > count || newPos < 1 {
		return domain.Checklist{}, fmt.Errorf(
			"position %d out of bound, max allowed is %d: %w",
			newPos, count, core_errors.ErrInvalidArgument,
		)
	}

	now := time.Now().UTC()

	if newPos < oldPosition {
		shiftQuery := `
			UPDATE taskana.checklists
			SET position = position + 1, version = version + 1, updated_at = $1
			WHERE user_id = $2
				AND task_id IS NOT DISTINCT FROM $3
				AND position >= $4
				AND position < $5
		`
		_, err = tx.Exec(
			ctx, shiftQuery, now, userID, checklist.TaskID, newPos, oldPosition,
		)
	} else if newPos > oldPosition {
		shiftQuery := `
			UPDATE taskana.checklists
			SET position = position - 1, version = version + 1, updated_at = $1
			WHERE user_id = $2
				AND task_id IS NOT DISTINCT FROM $3
				AND position > $4
				AND position <= $5
		`
		_, err = tx.Exec(
			ctx, shiftQuery, now, userID, checklist.TaskID, oldPosition, newPos,
		)
	}
	if err != nil {
		return domain.Checklist{},
			fmt.Errorf("failed to shift neighboring position: %w", err)
	}

	updateQuery := `
		UPDATE taskana.checklists
		SET 
		    position = $1,
		    updated_at = $2,
			version = version + 1
		WHERE id = $3 AND user_id = $4 AND version = $5
		RETURNING id, version, user_id, task_id, title, is_completed, position, 
		       created_at, updated_at
	`

	row := tx.QueryRow(
		ctx, updateQuery,
		newPos,
		now,
		checklist.ID,
		userID,
		checklist.Version,
	)

	checklistModel, err := scanChecklist(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Checklist{}, fmt.Errorf(
				"checklist with id='%s' concurrently accessed: %w",
				checklist.ID, core_errors.ErrConflict,
			)
		}
		return domain.Checklist{}, fmt.Errorf(
			"change position repository: %w", err,
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Checklist{}, fmt.Errorf(
			"failed to commit reordering transaction: %w", err,
		)
	}

	return checklistDomainFromModel(checklistModel), nil
}
