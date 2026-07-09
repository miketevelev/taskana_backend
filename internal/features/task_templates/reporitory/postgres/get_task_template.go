package task_templates_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
	core_postgres_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool"
)

func (r *TaskTemplateRepository) GetTaskTemplate(
	ctx context.Context,
	userID uuid.UUID,
	templateID uuid.UUID,
) (domain.TaskTemplate, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, user_id, project_id, heading_id, title, notes,
          recurrence_rule, recurrence_type, target_bucket, next_execution_date, 
          is_time_tracked, estimated_pomodoros, created_at, updated_at 
		FROM taskana.task_templates
		WHERE id = $1 AND user_id = $2
	`

	row := r.pool.QueryRow(ctx, query, templateID, userID)

	taskTemplateModel, err := scanTaskTemplate(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.TaskTemplate{}, fmt.Errorf(
				"task template with id %s not found: %w",
				templateID,
				core_errors.ErrNotFound,
			)
		}

		return domain.TaskTemplate{}, fmt.Errorf("scan error %w", err)
	}

	taskTemplateDomain := taskTemplateDomainFromModel(taskTemplateModel)

	return taskTemplateDomain, nil
}
