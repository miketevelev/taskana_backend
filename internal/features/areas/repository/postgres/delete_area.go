package areas_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
	core_postgres_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool"
)

func (r *AreasRepository) DeleteArea(
	ctx context.Context,
	userID uuid.UUID,
	areaID uuid.UUID,
) error {
	// todo: area must deleted all projects and tasks inside
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var pos int
	err = tx.QueryRow(
		ctx,
		"SELECT position FROM taskana.areas WHERE id = $1 AND user_id = $2",
		areaID, userID,
	).Scan(&pos)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return fmt.Errorf("area not found: %w", core_errors.ErrNotFound)
		}
		return fmt.Errorf("failed to get position: %w", err)
	}

	shiftQuery := `
		UPDATE taskana.areas 
		SET position = position - 1, updated_at = $1
		WHERE user_id = $2 AND position > $3
	`
	_, err = tx.Exec(ctx, shiftQuery, time.Now().UTC(), userID, pos)
	if err != nil {
		return fmt.Errorf("failed to shift positions: %w", err)
	}

	deleteQuery := `DELETE FROM taskana.areas WHERE id = $1 AND user_id = $2`
	cmdTag, err := tx.Exec(ctx, deleteQuery, areaID, userID)
	if err != nil {
		return fmt.Errorf("exec delete query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("area not found: %w", core_errors.ErrNotFound)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
