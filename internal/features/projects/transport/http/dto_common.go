package projects_transport_http

import (
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain/project"
)

type ProjectDTOResponse struct {
	ID          uuid.UUID                    `json:"id"`
	Version     int                          `json:"version"`
	UserID      uuid.UUID                    `json:"user_id"`
	AreaID      *uuid.UUID                   `json:"area_id,omitempty"`
	Title       string                       `json:"title"`
	Notes       *string                      `json:"notes,omitempty"`
	Status      domain_project.ProjectStatus `json:"status"`
	Position    int                          `json:"position"`
	Deadline    *time.Time                   `json:"deadline,omitempty"`
	CompletedAt *time.Time                   `json:"completed_at,omitempty"`
	CreatedAt   time.Time                    `json:"created_at,omitempty"`
	UpdatedAt   time.Time                    `json:"updated_at,omitempty"`
}

func projectDTOFromDomain(project domain_project.Project) ProjectDTOResponse {
	return ProjectDTOResponse{
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

func projectDTOsFromDomains(projects []domain_project.Project) []ProjectDTOResponse {
	result := make([]ProjectDTOResponse, len(projects))
	for i, project := range projects {
		result[i] = projectDTOFromDomain(project)
	}
	return result
}
