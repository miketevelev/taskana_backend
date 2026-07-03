package projects_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *ProjectService) CreateProject(
	ctx context.Context,
	userID uuid.UUID,
	project domain.Project,
) (domain.Project, error) {
	if project.ID == uuid.Nil {
		project.ID = uuid.New()
	}

	if project.Version == -1 {
		project.Version = 1
	}

	if err := project.Validate(); err != nil {
		return domain.Project{}, fmt.Errorf(
			"project validation failed: %w", err,
		)
	}

	createdProject, err := s.projectsRepository.CreateProject(
		ctx,
		userID,
		project,
	)
	if err != nil {
		return domain.Project{}, fmt.Errorf("failed to create project: %w", err)
	}

	return createdProject, nil
}
