package areas_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain/area"
)

func (r *AreasRepository) GetAreas(
	ctx context.Context,
	userID uuid.UUID,
	limit *int,
	offset *int,
) ([]domain_area.Area, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, user_id, title, position, created_at, updated_at
		FROM taskana.areas
		WHERE user_id = $1
		ORDER BY position ASC, created_at ASC
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.pool.Query(
		ctx,
		query,
		userID,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("select areas: %w", err)
	}
	defer rows.Close()

	var areaModels []AreaModel

	for rows.Next() {
		areaModel, err := scanArea(rows)
		if err != nil {
			return nil, fmt.Errorf("scan area from db: %w", err)
		}
		areaModels = append(areaModels, areaModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	areaDomains := areaDomainsFromModels(areaModels)

	return areaDomains, nil
}
