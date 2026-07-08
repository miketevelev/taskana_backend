package task_templates_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
)

func (r *TaskTemplateRepository) DeleteTaskTemplate(
	ctx context.Context,
	userID uuid.UUID,
	taskTemplateID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tag, err := r.pool.Exec(
		ctx,
		`DELETE FROM taskana.task_templates WHERE id = $1 AND user_id = $2`,
		taskTemplateID, userID,
	)
	if err != nil {
		return fmt.Errorf("delete task template: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return core_errors.ErrNotFound
	}

	return nil
}
