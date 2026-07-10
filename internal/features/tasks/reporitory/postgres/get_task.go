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

func (r *TasksRepository) GetTask(
	ctx context.Context,
	userID uuid.UUID,
	taskID uuid.UUID,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, user_id, project_id, heading_id, template_id, 
			title, notes, status, bucket, start_date, deadline, position, 
			is_time_tracked, estimated_pomodoros, completed_at, created_at, 
		    updated_at
		FROM taskana.tasks
		WHERE id = $1 AND user_id = $2
	`

	row := r.pool.QueryRow(ctx, query, taskID, userID)

	taskModel, err := scanTask(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf(
				"task with id %s not found: %w",
				taskID,
				core_errors.ErrNotFound,
			)
		}

		return domain.Task{}, fmt.Errorf("scan error %w", err)
	}

	taskDomain := taskDomainFromModel(taskModel)

	return taskDomain, nil
}
