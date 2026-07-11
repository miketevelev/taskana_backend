package checklists_postgres_repository

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

func (r *ChecklistsRepository) PatchChecklist(
	ctx context.Context,
	userID uuid.UUID,
	checklist domain.Checklist,
) (domain.Checklist, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE taskana.checklists
		SET 
			title = $1,
			is_completed = $2,
			updated_at = $3,
			version = version + 1
		WHERE id = $4 AND user_id = $5 AND version = $6
		RETURNING id, version, user_id, task_id, title, is_completed, position, 
		       created_at, updated_at
	`

	oldVersion := checklist.Version

	row := r.pool.QueryRow(
		ctx, query,
		checklist.Title,
		checklist.IsCompleted,
		time.Now().UTC(),
		checklist.ID,
		userID,
		oldVersion,
	)

	checklistModel, err := scanChecklist(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Checklist{}, fmt.Errorf(
				"checklist with id='%s' concurrently accessed: %w",
				checklist.ID,
				core_errors.ErrConflict,
			)
		}
		return domain.Checklist{}, fmt.Errorf(
			"patch checklist repository: %w", err,
		)
	}

	return checklistDomainFromModel(checklistModel), nil
}
