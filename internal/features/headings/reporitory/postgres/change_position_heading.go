package heading_postgres_repository

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

func (r *HeadingRepository) ChangePosition(
	ctx context.Context,
	userID uuid.UUID,
	heading domain.Heading,
	oldPosition int,
) (domain.Heading, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Heading{}, err
	}
	defer tx.Rollback(ctx)

	var count int
	countQuery := `
		SELECT COUNT(*) 
		FROM taskana.headings 
		WHERE user_id = $1 AND project_id = $2
	`
	if err := tx.QueryRow(
		ctx, countQuery, userID, heading.ProjectID,
	).Scan(&count); err != nil {
		return domain.Heading{}, fmt.Errorf(
			"failed to get headings count: %w", err,
		)
	}

	newPos := heading.Position
	if newPos > count || newPos < 1 {
		return domain.Heading{}, fmt.Errorf(
			"position %d out of bounds, max allowed is %d: %w",
			newPos, count, core_errors.ErrInvalidArgument,
		)
	}

	now := time.Now().UTC()

	if newPos < oldPosition {
		shiftQuery := `
			UPDATE taskana.headings 
			SET position = position + 1, version = version + 1, updated_at = $1
			WHERE user_id = $2 
			  AND project_id = $3 
			  AND position >= $4 
			  AND position < $5
		`
		_, err = tx.Exec(
			ctx, shiftQuery, now, userID, heading.ProjectID, newPos,
			oldPosition,
		)
	} else if newPos > oldPosition {
		shiftQuery := `
			UPDATE taskana.headings 
			SET position = position - 1, version = version + 1, updated_at = $1
			WHERE user_id = $2 
			  AND project_id = $3 
			  AND position > $4 
			  AND position <= $5
		`
		_, err = tx.Exec(
			ctx, shiftQuery, now, userID, heading.ProjectID, oldPosition,
			newPos,
		)
	}

	if err != nil {
		return domain.Heading{}, fmt.Errorf(
			"failed to shift neighboring positions: %w", err,
		)
	}

	updateQuery := `
		UPDATE taskana.headings
		SET 
			position = $1,
			updated_at = $2,
			version = version + 1
		WHERE id = $3 AND user_id = $4 AND version = $5
		RETURNING id, version, user_id, project_id, title, 
		    position, created_at, updated_at
	`
	row := tx.QueryRow(
		ctx, updateQuery,
		newPos,
		now,
		heading.ID,
		userID,
		heading.Version,
	)

	headingModel, err := scanHeading(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Heading{}, fmt.Errorf(
				"heading with id='%s' concurrently accessed: %w",
				heading.ID, core_errors.ErrConflict,
			)
		}
		return domain.Heading{}, fmt.Errorf(
			"change position repository: %w", err,
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Heading{}, fmt.Errorf(
			"failed to commit reordering transaction: %w", err,
		)
	}

	return headingDomainFromModel(headingModel), nil
}
