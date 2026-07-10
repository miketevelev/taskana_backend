package checklists_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (r *ChecklistsRepository) GetChecklists(
	ctx context.Context,
	userID uuid.UUID,
	limit *int,
	offset *int,
) ([]domain.Checklist, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT 
			c.id, c.version, c.user_id, c.task_id, c.title, 
			c.is_completed, c.position, c.created_at, c.updated_at
		FROM taskana.checklists c
		LEFT JOIN taskana.tasks a ON c.task_id = a.id
		WHERE c.user_id = $1
		ORDER BY 
			a.position ASC NULLS FIRST, 
			c.position ASC, 
			c.created_at ASC  
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.pool.Query(
		ctx, query,
		userID,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("select checklists: %w", err)
	}
	defer rows.Close()

	var checklistModels []ChecklistModel

	for rows.Next() {
		checklistModel, err := scanChecklist(rows)
		if err != nil {
			return nil, fmt.Errorf("scan checklist from db: %w", err)
		}
		checklistModels = append(checklistModels, checklistModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	checklistDomains := checklistDomainsFromModels(checklistModels)

	return checklistDomains, nil
}
