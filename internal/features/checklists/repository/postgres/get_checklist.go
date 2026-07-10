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

func (r *ChecklistsRepository) GetChecklist(
	ctx context.Context,
	userID uuid.UUID,
	checklistID uuid.UUID,
) (domain.Checklist, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, user_id, task_id, title, is_completed, position, 
		       created_at, updated_at
		FROM taskana.checklists
		WHERE id = $1 AND user_id = $2
	`

	row := r.pool.QueryRow(ctx, query, checklistID, userID)

	checklistModel, err := scanChecklist(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Checklist{}, fmt.Errorf(
				"checklist with id %s not found: %w",
				checklistID,
				core_errors.ErrNotFound,
			)
		}

		return domain.Checklist{}, fmt.Errorf("scan error %w", err)
	}

	checklistDomain := checklistDomainFromModel(checklistModel)

	return checklistDomain, nil
}
