package projects_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (r *ProjectRepository) GetProjects(
	ctx context.Context,
	userID uuid.UUID,
	limit *int,
	offset *int,
) ([]domain.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT 
			p.id, p.version, p.user_id, p.area_id, p.title, 
			p.notes, p.status, p.position, p.deadline, p.completed_at, p.created_at, p.updated_at
		FROM taskana.projects p
		LEFT JOIN taskana.areas a ON p.area_id = a.id
		WHERE p.user_id = $1
		ORDER BY 
			a.position ASC NULLS FIRST, 
			p.position ASC, 
			p.created_at ASC  
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.pool.Query(
		ctx,
		query,
		userID,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("select projects: %w", err)
	}
	defer rows.Close()

	var projectModels []ProjectModel

	for rows.Next() {
		projectModel, err := scanProject(rows)
		if err != nil {
			return nil, fmt.Errorf("scan projects from db: %w", err)
		}
		projectModels = append(projectModels, projectModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	projectDomains := projectDomainsFromModels(projectModels)

	return projectDomains, nil
}
