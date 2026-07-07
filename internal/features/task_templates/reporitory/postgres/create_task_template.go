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

func (r *TaskTemplateRepository) CreateTaskTemplate(
	ctx context.Context,
	userID uuid.UUID,
	taskTemplate domain.TaskTemplate,
) (domain.TaskTemplate, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO taskana.task_templates (
			id, user_id, project_id, heading_id, title, notes,
			recurrence_rule, recurrence_type, target_bucket, next_execution_date,
			is_time_tracked, estimated_pomodoros, created_at, updated_at)
		VALUES (
          $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
       )
		RETURNING id, version, user_id, project_id, heading_id, title, notes,
          recurrence_rule, recurrence_type, target_bucket, next_execution_date, 
          is_time_tracked, estimated_pomodoros, created_at, updated_at
	`

	row := r.pool.QueryRow(
		ctx, query,
		taskTemplate.ID,
		userID,
		taskTemplate.ProjectID,
		taskTemplate.HeadingID,
		taskTemplate.Title,
		taskTemplate.Notes,
		taskTemplate.RecurrenceRule,
		taskTemplate.RecurrenceType,
		taskTemplate.TargetBucket,
		taskTemplate.NextExecutionDate,
		taskTemplate.IsTimeTracked,
		taskTemplate.EstimatedPomodoros,
		taskTemplate.CreatedAt,
		taskTemplate.UpdatedAt,
	)

	taskTemplateModel, err := scanTaskTemplate(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolateForeignKey) {
			return domain.TaskTemplate{}, fmt.Errorf(
				"user not found for new task template: %w",
				core_errors.ErrNotFound,
			)
		}

		return domain.TaskTemplate{},
			fmt.Errorf("scan task template from db: %w", err)
	}

	taskTemplateDomain := taskTemplateDomainFromModel(taskTemplateModel)

	return taskTemplateDomain, nil
}
