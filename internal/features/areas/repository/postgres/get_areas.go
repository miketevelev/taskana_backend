package areas_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (r *AreasRepository) GetAreas(
	ctx context.Context,
	userID uuid.UUID,
	limit *int,
	offset *int,
) ([]domain.Area, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, user_id, title, position, created_at, updated_at
		FROM taskana.areas
		%s
		ORDER BY id ASC
		LIMIT $1 OFFSET $2;`

	args := []any{limit, offset}

	if userID != uuid.Nil {
		query = fmt.Sprintf(query, "WHERE user_id = $3")
		args = append(args, userID)
	} else {
		query = fmt.Sprintf(query, "")
	}

	rows, err := r.pool.Query(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf("select areas: %w", err)
	}
	defer rows.Close()

	var areaModels []AreaModel
	for rows.Next() {
		areaModel, err := scanArea(rows)
		if err != nil {
			return []domain.Area{}, fmt.Errorf("scan area from db: %w", err)
		}
		areaModels = append(areaModels, areaModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	areaDomains := areaDomainsFromModels(areaModels)

	return areaDomains, nil
}
