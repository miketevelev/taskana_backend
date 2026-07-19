package projects_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain/project"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
	core_postgres_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool"
)

func (r *ProjectRepository) ChangePosition(
	ctx context.Context,
	userID uuid.UUID,
	project domain_project.Project,
	oldPosition int,
) (domain_project.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain_project.Project{}, err
	}
	defer tx.Rollback(ctx)

	var count int
	countQuery := `
		SELECT COUNT(*) 
		FROM taskana.projects 
		WHERE user_id = $1 AND area_id IS NOT DISTINCT FROM $2
	`
	if err := tx.QueryRow(
		ctx, countQuery, userID, project.AreaID,
	).Scan(&count); err != nil {
		return domain_project.Project{}, fmt.Errorf(
			"failed to get projects count: %w", err,
		)
	}

	newPos := project.Position
	if newPos > count || newPos < 1 {
		return domain_project.Project{}, fmt.Errorf(
			"position %d out of bounds, max allowed is %d: %w",
			newPos, count, core_errors.ErrInvalidArgument,
		)
	}

	now := time.Now().UTC()

	if newPos < oldPosition {
		shiftQuery := `
			UPDATE taskana.projects 
			SET position = position + 1, version = version + 1, updated_at = $1
			WHERE user_id = $2 
			  AND area_id IS NOT DISTINCT FROM $3 
			  AND position >= $4 
			  AND position < $5
		`
		_, err = tx.Exec(
			ctx, shiftQuery, now, userID, project.AreaID, newPos, oldPosition,
		)
	} else if newPos > oldPosition {
		shiftQuery := `
			UPDATE taskana.projects 
			SET position = position - 1, version = version + 1, updated_at = $1
			WHERE user_id = $2 
			  AND area_id IS NOT DISTINCT FROM $3 
			  AND position > $4 
			  AND position <= $5
		`
		_, err = tx.Exec(
			ctx, shiftQuery, now, userID, project.AreaID, oldPosition, newPos,
		)
	}

	if err != nil {
		return domain_project.Project{}, fmt.Errorf(
			"failed to shift neighboring positions: %w", err,
		)
	}

	updateQuery := `
		UPDATE taskana.projects
		SET 
			position = $1,
			updated_at = $2,
			version = version + 1
		WHERE id = $3 AND user_id = $4 AND version = $5
		RETURNING id, version, user_id, area_id, title, notes, status, 
		    position, deadline, completed_at, created_at, updated_at
	`
	row := tx.QueryRow(
		ctx, updateQuery,
		newPos,
		now,
		project.ID,
		userID,
		project.Version,
	)

	projectModel, err := scanProject(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain_project.Project{}, fmt.Errorf(
				"project with id='%s' concurrently accessed: %w",
				project.ID, core_errors.ErrConflict,
			)
		}
		return domain_project.Project{}, fmt.Errorf(
			"change position repository: %w", err,
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain_project.Project{}, fmt.Errorf(
			"failed to commit reordering transaction: %w", err,
		)
	}

	return projectDomainFromModel(projectModel), nil
}
