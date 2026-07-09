package tasks_postgres_repository

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

func (r *TaskRepository) ChangePosition(
	ctx context.Context,
	userID uuid.UUID,
	task domain.Task,
	oldPosition int,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Task{}, fmt.Errorf(
			"failed to begin transaction: %w", err,
		)
	}
	defer tx.Rollback(ctx)

	var count int
	countQuery := `
		SELECT COUNT(*) 
		FROM taskana.tasks 
		WHERE user_id = $1 AND project_id IS NOT DISTINCT FROM $2
	`
	if err := tx.QueryRow(
		ctx, countQuery, userID, task.ProjectID,
	).Scan(&count); err != nil {
		return domain.Task{}, fmt.Errorf(
			"failed to get tasks count: %w", err,
		)
	}

	newPos := task.Position
	if newPos > count || newPos < 1 {
		return domain.Task{}, fmt.Errorf(
			"position %d out of bounds, max allowed is %d: %w",
			newPos, count, core_errors.ErrInvalidArgument,
		)
	}

	now := time.Now().UTC()

	if newPos < oldPosition {
		shiftQuery := `
			UPDATE taskana.tasks 
			SET position = position + 1, version = version + 1, updated_at = $1
			WHERE user_id = $2 
			  AND project_id IS NOT DISTINCT FROM $3 
			  AND position >= $4 
			  AND position < $5
		`
		_, err = tx.Exec(
			ctx, shiftQuery, now, userID, task.ProjectID, newPos, oldPosition,
		)
	} else if newPos > oldPosition {
		shiftQuery := `
			UPDATE taskana.tasks 
			SET position = position - 1, version = version + 1, updated_at = $1
			WHERE user_id = $2 
			  AND project_id IS NOT DISTINCT FROM $3 
			  AND position > $4 
			  AND position <= $5
		`
		_, err = tx.Exec(
			ctx, shiftQuery, now, userID, task.ProjectID, oldPosition, newPos,
		)
	}

	if err != nil {
		return domain.Task{}, fmt.Errorf(
			"failed to shift neighboring positions: %w", err,
		)
	}

	updateQuery := `
		UPDATE taskana.tasks
		SET 
			position = $1,
			updated_at = $2,
			version = version + 1
		WHERE id = $3 AND user_id = $4 AND version = $5
		RETURNING 
			id, version, user_id, project_id, heading_id, template_id, 
			title, notes, status, bucket, start_date, deadline, position, 
			is_time_tracked, estimated_pomodoros, completed_at, created_at, 
		    updated_at
	`
	row := tx.QueryRow(
		ctx, updateQuery,
		newPos,
		now,
		task.ID,
		userID,
		task.Version,
	)

	taskModel, err := scanTask(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf(
				"task with id='%s' concurrently accessed: %w",
				task.ID, core_errors.ErrConflict,
			)
		}
		return domain.Task{}, fmt.Errorf(
			"change position repository: %w", err,
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Task{}, fmt.Errorf(
			"failed to commit reordering transaction: %w", err,
		)
	}

	return taskDomainFromModel(taskModel), nil
}
