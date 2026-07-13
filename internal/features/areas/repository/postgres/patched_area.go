package areas_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain/area"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
	core_postgres_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool"
)

func (r *AreasRepository) PatchArea(
	ctx context.Context,
	userID uuid.UUID,
	area domain_area.Area,
) (domain_area.Area, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE taskana.areas
		SET 
			title = $1, 
			updated_at = $2,
			version = version + 1
		WHERE id = $3 AND user_id = $4 AND version = $5
		RETURNING id, version, user_id, title, position, created_at, updated_at
	`

	row := r.pool.QueryRow(
		ctx, query,
		area.Title,
		time.Now().UTC(),
		area.ID,
		userID,
		area.Version,
	)

	areaModel, err := scanArea(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain_area.Area{}, fmt.Errorf(
				"area concurrently accessed: %w", core_errors.ErrConflict,
			)
		}
		return domain_area.Area{}, fmt.Errorf("patch area repository: %w", err)
	}

	return areaDomainFromModel(areaModel), nil
}

func (r *AreasRepository) PatchAreaWithReordering(
	ctx context.Context,
	userID uuid.UUID,
	area domain_area.Area,
	oldPos int,
) (domain_area.Area, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain_area.Area{}, err
	}
	defer tx.Rollback(ctx)

	var count int
	countQuery := `SELECT COUNT(*) FROM taskana.areas WHERE user_id = $1`
	if err := tx.QueryRow(ctx, countQuery, userID).Scan(&count); err != nil {
		return domain_area.Area{}, fmt.Errorf(
			"failed to get areas count: %w", err,
		)
	}

	newPos := area.Position
	if newPos > count || newPos < 1 {
		return domain_area.Area{}, fmt.Errorf(
			"position %d out of bounds, max allowed is %d: %w",
			newPos, count, core_errors.ErrInvalidArgument,
		)
	}

	now := time.Now().UTC()

	if newPos < oldPos {
		shiftQuery := `
			UPDATE taskana.areas 
			SET position = position + 1, version = version + 1, updated_at = $1
			WHERE user_id = $2 AND position >= $3 AND position < $4
		`
		_, err = tx.Exec(ctx, shiftQuery, now, userID, newPos, oldPos)
	} else if newPos > oldPos {
		shiftQuery := `
			UPDATE taskana.areas 
			SET position = position - 1, version = version + 1, updated_at = $1
			WHERE user_id = $2 AND position > $3 AND position <= $4
		`
		_, err = tx.Exec(ctx, shiftQuery, now, userID, oldPos, newPos)
	}

	if err != nil {
		return domain_area.Area{}, fmt.Errorf(
			"failed to shift neighboring positions: %w", err,
		)
	}

	updateQuery := `
		UPDATE taskana.areas
		SET 
			title = $1, 
			updated_at = $2,
			position = $3,
			version = version + 1
		WHERE id = $4 AND user_id = $5 AND version = $6
		RETURNING id, version, user_id, title, position, created_at, updated_at
	`
	row := tx.QueryRow(
		ctx, updateQuery,
		area.Title,
		now,
		newPos,
		area.ID,
		userID,
		area.Version,
	)

	areaModel, err := scanArea(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain_area.Area{}, fmt.Errorf(
				"area with id='%s' concurrently accessed: %w",
				area.ID,
				core_errors.ErrConflict,
			)
		}
		return domain_area.Area{}, fmt.Errorf("patch area repository: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain_area.Area{}, fmt.Errorf(
			"failed to commit reordering transaction: %w", err,
		)
	}

	areaDomain := areaDomainFromModel(areaModel)

	return areaDomain, nil
}
