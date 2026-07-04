package projects_postgres_repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type ProjectModel struct {
	ID          uuid.UUID            `json:"id"`
	Version     int                  `json:"version"`
	UserID      uuid.UUID            `json:"user_id"`
	AreaID      *uuid.UUID           `json:"area_id,omitempty"`
	Title       string               `json:"title"`
	Notes       *string              `json:"notes,omitempty"`
	Status      domain.ProjectStatus `json:"status"`
	Position    int                  `json:"position"`
	Deadline    *time.Time           `json:"deadline,omitempty"`
	CompletedAt *time.Time           `json:"completed_at,omitempty"`
	CreatedAt   time.Time            `json:"created_at,omitempty"`
	UpdatedAt   time.Time            `json:"updated_at,omitempty"`
}

func projectDomainFromModel(projectModel ProjectModel) domain.Project {
	return domain.Project{
		ID:          projectModel.ID,
		Version:     projectModel.Version,
		UserID:      projectModel.UserID,
		AreaID:      projectModel.AreaID,
		Title:       projectModel.Title,
		Notes:       projectModel.Notes,
		Status:      projectModel.Status,
		Position:    projectModel.Position,
		Deadline:    projectModel.Deadline,
		CompletedAt: projectModel.CompletedAt,
		CreatedAt:   projectModel.CreatedAt,
		UpdatedAt:   projectModel.UpdatedAt,
	}
}

func projectDomainsFromModels(projects []ProjectModel) []domain.Project {
	if len(projects) == 0 {
		return []domain.Project{}
	}
	projectDomains := make([]domain.Project, len(projects))

	for i, project := range projects {
		projectDomains[i] = domain.Project{
			ID:          project.ID,
			Version:     project.Version,
			UserID:      project.UserID,
			AreaID:      project.AreaID,
			Title:       project.Title,
			Notes:       project.Notes,
			Status:      project.Status,
			Position:    project.Position,
			Deadline:    project.Deadline,
			CompletedAt: project.CompletedAt,
			CreatedAt:   project.CreatedAt,
			UpdatedAt:   project.UpdatedAt,
		}
	}

	return projectDomains
}

func scanProject(row interface{ Scan(dest ...any) error }) (
	ProjectModel,
	error,
) {
	var projectModel ProjectModel
	err := row.Scan(
		&projectModel.ID,
		&projectModel.Version,
		&projectModel.UserID,
		&projectModel.AreaID,
		&projectModel.Title,
		&projectModel.Notes,
		&projectModel.Status,
		&projectModel.Position,
		&projectModel.Deadline,
		&projectModel.CompletedAt,
		&projectModel.CreatedAt,
		&projectModel.UpdatedAt,
	)
	if err != nil {
		return ProjectModel{}, fmt.Errorf("scan project: %w", err)
	}
	return projectModel, nil
}
