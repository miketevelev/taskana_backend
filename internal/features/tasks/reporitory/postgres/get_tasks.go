package tasks_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (r *TasksRepository) GetTasks(
	ctx context.Context,
	userID uuid.UUID,
	limit *int,
	offset *int,
) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
       SELECT 
          t.id, t.version, t.user_id, t.project_id, t.heading_id, t.template_id, 
          t.title, t.notes, t.status, t.bucket, t.start_date, t.deadline, t.position, 
          t.is_time_tracked, t.estimated_pomodoros, t.completed_at, t.created_at, 
          t.updated_at
       FROM taskana.tasks t
       LEFT JOIN taskana.projects p ON t.project_id = p.id
       WHERE t.user_id = $1
       ORDER BY 
          p.position ASC NULLS FIRST, 
          t.position ASC, 
          t.created_at ASC  
       LIMIT $2 OFFSET $3;
    `

	rows, err := r.pool.Query(
		ctx,
		query,
		userID,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("select tasks: %w", err)
	}
	defer rows.Close()

	var taskModels []TaskModel

	for rows.Next() {
		taskModel, err := scanTask(rows)
		if err != nil {
			return nil, fmt.Errorf("scan tasks from db: %w", err)
		}
		taskModels = append(taskModels, taskModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	taskDomains := taskDomainsFromModels(taskModels)

	return taskDomains, nil
}
