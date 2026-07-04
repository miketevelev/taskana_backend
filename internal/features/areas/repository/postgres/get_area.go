package areas_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
	core_postgres_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool"
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
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Area{}, fmt.Errorf(
				"area with id %s not found: %w",
				areaID,
				core_errors.ErrNotFound,
			)
		}

		return domain.Area{}, fmt.Errorf("scan error %w", err)
	}

	areaDomain := areaDomainFromModel(areaModel)

	return areaDomain, nil
}
