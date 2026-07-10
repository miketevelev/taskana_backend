package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
	core_postgres_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool"
)

type TasksRepository struct {
	pool core_postgres_pool.Pool
}

func NewTasksRepository(
	pool core_postgres_pool.Pool,
) *TasksRepository {
	return &TasksRepository{
		pool: pool,
	}
}

func (r *TasksRepository) UpdateTemplateNextExecution(
	ctx context.Context,
	userID, templateID uuid.UUID,
	nextDate time.Time,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE taskana.task_templates
		SET next_execution_date = $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3
	`
	tag, err := r.pool.Exec(ctx, query, nextDate, templateID, userID)
	if err != nil {
		return fmt.Errorf("update template: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return core_errors.ErrNotFound
	}
	return nil
}

func (r *TasksRepository) ListDueFixedTemplates(
	ctx context.Context,
	asOf time.Time,
) ([]domain.TaskTemplate, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, user_id, project_id, heading_id, title, notes,
		       recurrence_rule, recurrence_type, target_bucket, next_execution_date,
		       is_time_tracked, estimated_pomodoros, created_at, updated_at
		FROM taskana.task_templates
		WHERE recurrence_type = 'fixed'
		  AND next_execution_date <= $1
		ORDER BY next_execution_date ASC
	`

	rows, err := r.pool.Query(ctx, query, asOf)
	if err != nil {
		return nil, fmt.Errorf("list due templates: %w", err)
	}
	defer rows.Close()

	var result []domain.TaskTemplate
	for rows.Next() {
		var m TaskTemplateModel
		err := rows.Scan(
			&m.ID,
			&m.Version,
			&m.UserID,
			&m.ProjectID,
			&m.HeadingID,
			&m.Title,
			&m.Notes,
			&m.RecurrenceRule,
			&m.RecurrenceType,
			&m.TargetBucket,
			&m.NextExecutionDate,
			&m.IsTimeTracked,
			&m.EstimatedPomodoros,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan template: %w", err)
		}
		result = append(result, taskTemplateDomainFromModel(m))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate templates: %w", err)
	}

	return result, nil
}

func (r *TasksRepository) ProcessNextDueFixedTemplateTx(
	ctx context.Context,
	asOf time.Time,
	processFn func(template domain.TaskTemplate) (
		domain.Task,
		time.Time,
		error,
	),
) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	queryLock := `
		SELECT id, version, user_id, project_id, heading_id, title, notes,
		       recurrence_rule, recurrence_type, target_bucket, next_execution_date,
		       is_time_tracked, estimated_pomodoros, created_at, updated_at
		FROM taskana.task_templates
		WHERE recurrence_type = 'fixed'
		  AND next_execution_date <= $1
		ORDER BY next_execution_date ASC
		FOR UPDATE SKIP LOCKED
		LIMIT 1
	`

	var m TaskTemplateModel
	err = tx.QueryRow(ctx, queryLock, asOf).Scan(
		&m.ID, &m.Version, &m.UserID, &m.ProjectID, &m.HeadingID, &m.Title,
		&m.Notes, &m.RecurrenceRule, &m.RecurrenceType, &m.TargetBucket,
		&m.NextExecutionDate, &m.IsTimeTracked, &m.EstimatedPomodoros,
		&m.CreatedAt, &m.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("lock next template: %w", err)
	}

	template := taskTemplateDomainFromModel(m)

	newTask, nextDate, err := processFn(template)
	if err != nil {
		return false, fmt.Errorf("process domain logic: %w", err)
	}

	insertTaskQuery := `
		INSERT INTO taskana.tasks (
			id, version, user_id, project_id, heading_id, template_id,
			title, notes, status, bucket, start_date, deadline,
			position, is_time_tracked, estimated_pomodoros, completed_at,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12,
			$13, $14, $15, $16, $17, $18
		)
	`

	_, err = tx.Exec(
		ctx,
		insertTaskQuery,
		newTask.ID,
		newTask.Version,
		newTask.UserID,
		newTask.ProjectID,
		newTask.HeadingID,
		newTask.TemplateID,
		newTask.Title,
		newTask.Notes,
		newTask.Status,
		newTask.Bucket,
		newTask.StartDate,
		newTask.Deadline,
		newTask.Position,
		newTask.IsTimeTracked,
		newTask.EstimatedPomodoros,
		newTask.CompletedAt,
		newTask.CreatedAt,
		newTask.UpdatedAt,
	)
	if err != nil {
		return false, fmt.Errorf("insert generated task: %w", err)
	}

	updateTemplateQuery := `
		UPDATE taskana.task_templates
		SET next_execution_date = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err = tx.Exec(ctx, updateTemplateQuery, nextDate, template.ID)
	if err != nil {
		return false, fmt.Errorf("update template next date: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return false, fmt.Errorf("commit tx: %w", err)
	}

	return true, nil
}
