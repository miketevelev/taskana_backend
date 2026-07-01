package areas_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (r *AreasRepository) GetArea(
	ctx context.Context,
	userID uuid.UUID,
	areaID uuid.UUID,
) (domain.Area, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, user_id, title, position, created_at, updated_at
		FROM taskana.areas
		WHERE id = $1 AND user_id = $2
	`

	row := r.pool.QueryRow(ctx, query, areaID, userID)

	areaModel, err := scanArea(row)
	if err != nil {
		return domain.Area{}, fmt.Errorf("scan area from db: %w", err)
	}

	areaDomain := areaDomainFromModel(areaModel)

	return areaDomain, nil
}
