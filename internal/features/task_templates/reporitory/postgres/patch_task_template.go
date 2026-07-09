package task_templates_postgres_repository

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

func (r *TaskTemplateRepository) PatchTaskTemplate(
	ctx context.Context,
	userID uuid.UUID,
	taskTemplate domain.TaskTemplate,
) (domain.TaskTemplate, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE taskana.task_templates
		SET project_id = $1,
		    heading_id = $2,
		    title = $3,
		    notes = $4,
		    recurrence_rule = $5,
		    recurrence_type = $6,
		    target_bucket = $7,
		    next_execution_date = $8,
		    is_time_tracked = $9,
		    estimated_pomodoros = $10,
		    updated_at = $11,
		    version = version + 1
		WHERE id = $12 AND user_id = $13 AND version = $14
		RETURNING id, version, user_id, project_id, heading_id, title, notes,
          recurrence_rule, recurrence_type, target_bucket, next_execution_date, 
          is_time_tracked, estimated_pomodoros, created_at, updated_at 
	`

	row := r.pool.QueryRow(
		ctx, query,
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
		time.Now().UTC(),
		taskTemplate.ID,
		userID,
		taskTemplate.Version,
	)

	taskTemplateModel, err := scanTaskTemplate(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.TaskTemplate{}, fmt.Errorf(
				"task template with id='%s' concurrently accessed: %w",
				taskTemplate.ID,
				core_errors.ErrConflict,
			)
		}
		return domain.TaskTemplate{}, fmt.Errorf(
			"patch task template repository: %w", err,
		)
	}

	return taskTemplateDomainFromModel(taskTemplateModel), nil
}
