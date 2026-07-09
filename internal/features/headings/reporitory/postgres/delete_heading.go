package heading_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
	core_postgres_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool"
)

func (r *HeadingRepository) DeleteHeading(
	ctx context.Context,
	userID uuid.UUID,
	headingID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var projectID uuid.UUID
	var position int

	err = tx.QueryRow(
		ctx, `
			SELECT project_id, position
			FROM taskana.headings
			WHERE id = $1 AND user_id = $2`,
		headingID, userID,
	).Scan(&projectID, &position)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return fmt.Errorf("heading not found: %w", core_errors.ErrNotFound)
		}
		return fmt.Errorf("failed to find heading to delete: %w", err)
	}

	shiftQuery := `
		UPDATE taskana.headings
		SET position = position - 1
		WHERE user_id = $1
			AND project_id = $2
			AND position > $3
	`
	if _, err = tx.Exec(
		ctx, shiftQuery, userID, projectID, position,
	); err != nil {
		return fmt.Errorf("failed to shift positions: %w", err)
	}

	deleteQuery := `DELETE FROM taskana.headings WHERE id = $1 AND user_id = $2`
	if _, err = tx.Exec(ctx, deleteQuery, headingID, userID); err != nil {
		return fmt.Errorf("failed to delete heading: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
