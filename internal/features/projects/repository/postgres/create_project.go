package projects_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
	core_postgres_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool"
)

func (r *ProjectRepository) CreateProject(
	ctx context.Context,
	userID uuid.UUID,
	project domain.Project,
) (domain.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO taskana.projects (id, version, user_id, area_id, title, 
notes, status, position, deadline, completed_at, created_at, updated_at)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, 
			(
				SELECT COALESCE(MAX(position), 0) + 1 
				FROM taskana.projects 
				WHERE user_id = $3 AND area_id IS NOT DISTINCT FROM $4
			),
			$8, $9, $10, $11
		)
		RETURNING id, version, user_id, area_id, title, 
notes, status, position, deadline, completed_at, created_at, updated_at
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		project.ID,
		project.Version,
		userID,
		project.AreaID,
		project.Title,
		project.Notes,
		project.Status,
		project.Deadline,
		project.CompletedAt,
		project.CreatedAt,
		project.UpdatedAt,
	)

	projectModel, err := scanProject(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolateForeignKey) {
			return domain.Project{}, fmt.Errorf(
				"user or area not found for new project: %w",
				core_errors.ErrNotFound,
			)
		}

		return domain.Project{}, fmt.Errorf("scan project from db: %w", err)
	}

	projectDomain := projectDomainFromModel(projectModel)

	return projectDomain, nil
}
