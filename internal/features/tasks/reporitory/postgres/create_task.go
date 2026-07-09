package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
	core_postgres_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool"
)

func (r *TaskRepository) CreateTask(
	ctx context.Context,
	userID uuid.UUID,
	task domain.Task,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO taskana.tasks (
			id, version, user_id, project_id, heading_id, template_id, 
			title, notes, status, bucket, start_date, deadline, position, 
			is_time_tracked, estimated_pomodoros, completed_at, created_at, updated_at
		) 
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
			(
				SELECT COALESCE(MAX(position), 0) + 1 
				FROM taskana.tasks 
				WHERE user_id = $3 AND project_id IS NOT DISTINCT FROM $4
			),
			$13, $14, $15, $16, $17
		)
		RETURNING 
			id, version, user_id, project_id, heading_id, template_id, 
			title, notes, status, bucket, start_date, deadline, position, 
			is_time_tracked, estimated_pomodoros, completed_at, created_at, updated_at
	`

	row := r.pool.QueryRow(
		ctx, query,
		task.ID,
		task.Version,
		userID,
		task.ProjectID,
		task.HeadingID,
		task.TemplateID,
		task.Title,
		task.Notes,
		task.Status,
		task.Bucket,
		task.StartDate,
		task.Deadline,
		task.IsTimeTracked,
		task.EstimatedPomodoros,
		task.CompletedAt,
		task.CreatedAt,
		task.UpdatedAt,
	)

	taskModel, err := scanTask(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolateForeignKey) {
			return domain.Task{}, fmt.Errorf(
				"user or project not found for new task: %w",
				core_errors.ErrNotFound,
			)
		}
		return domain.Task{}, fmt.Errorf("scan task from db: %w", err)
	}

	taskDomain := taskDomainFromModel(taskModel)

	return taskDomain, nil
}
