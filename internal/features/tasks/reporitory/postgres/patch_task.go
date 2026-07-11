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

func (r *TasksRepository) PatchTask(
	ctx context.Context,
	userID uuid.UUID,
	task domain.Task,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE taskana.tasks
		SET 
			project_id = $1,
			heading_id = $2,
			title = $3,
			notes = $4,
			status = $5,
			bucket = $6,
			start_date = $7,
			deadline = $8,
			position = $9,
			is_time_tracked = $10,
			estimated_pomodoros = $11,
			completed_at = $12,
			updated_at = $13,
			version = version + 1
		WHERE id = $14 AND user_id = $15 AND version = $16
		RETURNING id, version, user_id, project_id, heading_id, template_id, 
		          title, notes, status, bucket, start_date, deadline, position, 
		          is_time_tracked, estimated_pomodoros, completed_at, 
		    	  created_at, updated_at
	`

	oldVersion := task.Version

	row := r.pool.QueryRow(
		ctx, query,
		task.ProjectID,
		task.HeadingID,
		task.Title,
		task.Notes,
		task.Status,
		task.Bucket,
		task.StartDate,
		task.Deadline,
		task.Position,
		task.IsTimeTracked,
		task.EstimatedPomodoros,
		task.CompletedAt,
		time.Now().UTC(),
		task.ID,
		userID,
		oldVersion,
	)

	taskModel, err := scanTask(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf(
				"task with id='%s' concurrently accessed: %w",
				task.ID,
				core_errors.ErrConflict,
			)
		}
		return domain.Task{}, fmt.Errorf(
			"patch task repository: %w", err,
		)
	}

	return taskDomainFromModel(taskModel), nil
}

func (r *TasksRepository) PatchTaskWithProjectChange(
	ctx context.Context,
	userID uuid.UUID,
	task domain.Task,
	oldPosition int,
	oldProjectID *uuid.UUID,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Task{}, err
	}
	defer tx.Rollback(ctx)

	now := time.Now().UTC()

	var countNewProject int
	countNewQuery := `
		SELECT COUNT(*) 
		FROM taskana.tasks 
		WHERE user_id = $1 AND project_id IS NOT DISTINCT FROM $2`

	if err := tx.QueryRow(
		ctx, countNewQuery, userID, task.ProjectID,
	).Scan(&countNewProject); err != nil {
		return domain.Task{}, fmt.Errorf(
			"failed to get count in new project: %w", err,
		)
	}

	newPosition := countNewProject + 1

	shiftOldQuery := `
		UPDATE taskana.tasks 
		SET position = position - 1, version = version + 1, updated_at = $1
		WHERE user_id = $2 
		  AND project_id IS NOT DISTINCT FROM $3 
		  AND position > $4
	`
	if _, err = tx.Exec(
		ctx, shiftOldQuery, now, userID, oldProjectID, oldPosition,
	); err != nil {
		return domain.Task{}, fmt.Errorf(
			"failed to shift in old project: %w", err,
		)
	}

	updateQuery := `
		UPDATE taskana.tasks
		SET 
			project_id = $1, heading_id = $2, title = $3, notes = $4,
			status = $5, bucket = $6, start_date = $7, deadline = $8,
			position = $9, is_time_tracked = $10, estimated_pomodoros = $11,
			completed_at = $12, updated_at = $13, version = version + 1
		WHERE id = $14 AND user_id = $15 AND version = $16
		RETURNING id, version, user_id, project_id, heading_id, template_id, 
		          title, notes, status, bucket, start_date, deadline, position, 
		          is_time_tracked, estimated_pomodoros, completed_at, 
		          created_at, updated_at
	`
	row := tx.QueryRow(
		ctx, updateQuery,
		task.ProjectID, task.HeadingID, task.Title, task.Notes,
		task.Status, task.Bucket, task.StartDate, task.Deadline,
		newPosition,
		task.IsTimeTracked, task.EstimatedPomodoros, task.CompletedAt,
		now, task.ID, userID, task.Version,
	)

	taskModel, err := scanTask(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf(
				"task concurrently accessed: %w", core_errors.ErrConflict,
			)
		}
		return domain.Task{}, fmt.Errorf("patch task repository: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Task{}, fmt.Errorf(
			"failed to commit project change transaction: %w", err,
		)
	}

	return taskDomainFromModel(taskModel), nil
}
