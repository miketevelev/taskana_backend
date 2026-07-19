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

func (r *ProjectRepository) PatchProject(
	ctx context.Context,
	userID uuid.UUID,
	project domain_project.Project,
) (domain_project.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE taskana.projects
		SET 
			area_id = $1,
			title = $2,
			notes = $3,
			status = $4,
			position = $5,
			deadline = $6,
			completed_at = $7,
			updated_at = $8,
			version = version + 1
		WHERE id = $9 AND user_id = $10 AND version = $11
		RETURNING id, version, user_id, area_id, title, notes, status, 
		    position, deadline, completed_at, created_at, updated_at
	`

	oldVersion := project.Version

	row := r.pool.QueryRow(
		ctx, query,
		project.AreaID,
		project.Title,
		project.Notes,
		project.Status,
		project.Position,
		project.Deadline,
		project.CompletedAt,
		time.Now().UTC(),
		project.ID,
		userID,
		oldVersion,
	)

	projectModel, err := scanProject(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain_project.Project{}, fmt.Errorf(
				"project with id='%s' concurrently accessed: %w",
				project.ID,
				core_errors.ErrConflict,
			)
		}
		return domain_project.Project{}, fmt.Errorf(
			"patch project repository: %w", err,
		)
	}

	return projectDomainFromModel(projectModel), nil
}

func (r *ProjectRepository) PatchProjectWithAreaChange(
	ctx context.Context,
	userID uuid.UUID,
	project domain_project.Project,
	oldPosition int,
	oldAreaID *uuid.UUID,
) (domain_project.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain_project.Project{}, err
	}
	defer tx.Rollback(ctx)

	now := time.Now().UTC()

	var countNewArea int
	countNewQuery := `
		SELECT COUNT(*) 
		FROM taskana.projects 
		WHERE user_id = $1 AND area_id IS NOT DISTINCT FROM $2`

	if err := tx.QueryRow(
		ctx, countNewQuery, userID, project.AreaID,
	).Scan(&countNewArea); err != nil {
		return domain_project.Project{}, fmt.Errorf(
			"failed to get count in new area: %w", err,
		)
	}

	newPosition := countNewArea + 1

	shiftOldQuery := `
		UPDATE taskana.projects 
		SET position = position - 1, version = version + 1, updated_at = $1
		WHERE user_id = $2 
		  AND area_id IS NOT DISTINCT FROM $3 
		  AND position > $4
	`
	if _, err = tx.Exec(
		ctx, shiftOldQuery, now, userID, oldAreaID, oldPosition,
	); err != nil {
		return domain_project.Project{}, fmt.Errorf(
			"failed to shift in old area: %w", err,
		)
	}

	updateQuery := `
		UPDATE taskana.projects
		SET 
			area_id = $1, title = $2, notes = $3, status = $4,
			position = $5, deadline = $6, completed_at = $7,
			updated_at = $8, version = version + 1
		WHERE id = $9 AND user_id = $10 AND version = $11
		RETURNING id, version, user_id, area_id, title, notes, status, 
		    position, deadline, completed_at, created_at, updated_at
	`
	row := tx.QueryRow(
		ctx, updateQuery,
		project.AreaID, project.Title, project.Notes, project.Status,
		newPosition,
		project.Deadline, project.CompletedAt,
		now, project.ID, userID, project.Version,
	)

	projectModel, err := scanProject(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain_project.Project{}, fmt.Errorf(
				"project concurrently accessed: %w", core_errors.ErrConflict,
			)
		}
		return domain_project.Project{}, fmt.Errorf(
			"patch project repository: %w", err,
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain_project.Project{}, fmt.Errorf(
			"failed to commit area change transaction: %w", err,
		)
	}

	return projectDomainFromModel(projectModel), nil
}
