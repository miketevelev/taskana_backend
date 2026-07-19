package projects_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain/project"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
	core_postgres_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool"
)

func (r *ProjectRepository) GetProject(
	ctx context.Context,
	userID uuid.UUID,
	projectID uuid.UUID,
) (domain_project.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, user_id, area_id, title, 
notes, status, position, deadline, completed_at, created_at, updated_at
		FROM taskana.projects
		WHERE id = $1 AND user_id = $2
	`

	row := r.pool.QueryRow(ctx, query, projectID, userID)

	projectModel, err := scanProject(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain_project.Project{}, fmt.Errorf(
				"project with id %s not found: %w",
				projectID,
				core_errors.ErrNotFound,
			)
		}

		return domain_project.Project{}, fmt.Errorf("scan error %w", err)
	}

	projectDomain := projectDomainFromModel(projectModel)

	return projectDomain, nil
}
