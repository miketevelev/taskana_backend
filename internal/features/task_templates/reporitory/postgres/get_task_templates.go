package task_templates_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (r *TaskTemplateRepository) GetTaskTemplates(
	ctx context.Context,
	userId uuid.UUID,
	limit *int,
	offset *int,
) ([]domain.TaskTemplate, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, user_id, project_id, heading_id, title, notes,
          recurrence_rule, recurrence_type, target_bucket, next_execution_date, 
          is_time_tracked, estimated_pomodoros, created_at, updated_at 
		FROM taskana.task_templates
		WHERE user_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.pool.Query(ctx, query, userId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select task templates: %w", err)
	}
	defer rows.Close()

	var taskTemplatesModels []TaskTemplateModel

	for rows.Next() {
		taskTemplateModel, err := scanTaskTemplate(rows)
		if err != nil {
			return nil, fmt.Errorf("scan task templates from db: %w", err)
		}
		taskTemplatesModels = append(taskTemplatesModels, taskTemplateModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	taskTemplatesDomains := taskTemplateDomainsFromModels(taskTemplatesModels)

	return taskTemplatesDomains, nil
}
