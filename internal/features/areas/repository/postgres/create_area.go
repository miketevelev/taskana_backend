package areas_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (r *AreasRepository) CreateArea(
	ctx context.Context,
	userID uuid.UUID,
	area domain.Area,
) (domain.Area, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO taskana.areas (id, user_id, title, position, created_at, updated_at)
		VALUES (
			$1, $2, $3, 
			(
				SELECT COALESCE(MAX(position), 0) + 1 
				FROM taskana.areas 
				WHERE user_id = $2 
			),
			$4, $5
		)
		RETURNING id, version, user_id, title, position, created_at, updated_at
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		area.ID,
		area.UserID,
		area.Title,
		area.CreatedAt,
		area.UpdatedAt,
	)

	var areaModel AreaModel
	err := row.Scan(
		&areaModel.ID,
		&areaModel.Version,
		&areaModel.UserID,
		&areaModel.Title,
		&areaModel.Position,
		&areaModel.CreatedAt,
		&areaModel.UpdatedAt,
	)
	if err != nil {
		return domain.Area{}, fmt.Errorf("scan area from db: %w", err)
	}

	areaDomain := areaDomainFromModel(areaModel)

	return areaDomain, nil
}
