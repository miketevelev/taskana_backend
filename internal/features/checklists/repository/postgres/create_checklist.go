package checklists_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
	core_postgres_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool"
)

func (r *ChecklistsRepository) CreateChecklist(
	ctx context.Context,
	userID uuid.UUID,
	checklist domain.Checklist,
) (domain.Checklist, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO taskana.checklists (id, version, user_id, task_id, title, 
		                                is_completed, position, created_at, 
		                                updated_at)
		VALUES (
			$1, $2, $3, $4, $5, $6, 
			(
				SELECT COALESCE(MAX(position), 0) + 1 
				FROM taskana.checklists
				WHERE user_id = $3 AND task_id IS NOT DISTINCT FROM $4
			),
			$7, $8
		)
		RETURNING id, version, user_id, task_id, title, is_completed, position, 
		    created_at, updated_at
	`

	row := r.pool.QueryRow(
		ctx, query,
		checklist.ID,
		checklist.Version,
		userID,
		checklist.TaskID,
		checklist.Title,
		checklist.IsCompleted,
		checklist.CreatedAt,
		checklist.UpdatedAt,
	)

	checklistModel, err := scanChecklist(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolateForeignKey) {
			return domain.Checklist{}, fmt.Errorf(
				"user or task not found for new checklist: %w",
				core_errors.ErrNotFound,
			)
		}

		return domain.Checklist{}, fmt.Errorf("scan checklist from db: %w", err)
	}

	checklistDomain := checklistDomainFromModel(checklistModel)

	return checklistDomain, nil
}
