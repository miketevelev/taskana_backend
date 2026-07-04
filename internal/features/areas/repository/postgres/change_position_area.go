package areas_postgres_repository

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

func (r *AreasRepository) ChangePosition(
	ctx context.Context,
	userID uuid.UUID,
	area domain.Area,
	oldPosition int,
) (domain.Area, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Area{}, err
	}
	defer tx.Rollback(ctx)

	var count int
	if err := tx.QueryRow(
		ctx, `SELECT COUNT(*) FROM taskana.areas WHERE user_id = $1`, userID,
	).Scan(&count); err != nil {
		return domain.Area{}, fmt.Errorf(
			"failed to get areas count: %w", err,
		)
	}

	newPos := area.Position
	if newPos > count || newPos < 1 {
		return domain.Area{}, fmt.Errorf(
			"position %d out of bounds, max allowed is %d: %w",
			newPos, count, core_errors.ErrInvalidArgument,
		)
	}

	now := time.Now().UTC()

	if newPos < oldPosition {
		shiftQuery := `
			UPDATE taskana.areas 
			SET position = position + 1, version = version + 1, updated_at = $1
			WHERE user_id = $2 AND position >= $3 AND position < $4
		`
		_, err = tx.Exec(ctx, shiftQuery, now, userID, newPos, oldPosition)
	} else if newPos > oldPosition {
		shiftQuery := `
			UPDATE taskana.areas 
			SET position = position - 1, version = version + 1, updated_at = $1
			WHERE user_id = $2 AND position > $3 AND position <= $4
		`
		_, err = tx.Exec(ctx, shiftQuery, now, userID, oldPosition, newPos)
	}

	if err != nil {
		return domain.Area{}, fmt.Errorf(
			"failed to shift neighboring positions: %w", err,
		)
	}

	updateQuery := `
		UPDATE taskana.areas
		SET 
			position = $1,
			updated_at = $2,
			version = version + 1
		WHERE id = $3 AND user_id = $4 AND version = $5
		RETURNING id, version, user_id, title, position, created_at, updated_at
	`
	row := tx.QueryRow(
		ctx, updateQuery,
		newPos,
		now,
		area.ID,
		userID,
		area.Version,
	)

	areaModel, err := scanArea(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Area{}, fmt.Errorf(
				"area concurrently accessed: %w", core_errors.ErrConflict,
			)
		}
		return domain.Area{}, fmt.Errorf("change position repository: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Area{}, fmt.Errorf(
			"failed to commit reordering transaction: %w", err,
		)
	}

	return areaDomainFromModel(areaModel), nil
}
