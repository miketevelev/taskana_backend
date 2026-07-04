package projects_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *ProjectService) ChangePosition(
	ctx context.Context,
	userID uuid.UUID,
	projectID uuid.UUID,
	newPosition int,
) (domain.Project, error) {
	project, err := s.projectsRepository.GetProject(ctx, userID, projectID)
	if err != nil {
		return domain.Project{}, fmt.Errorf(
			"error getting project: %w", err,
		)
	}

	oldPosition := project.Position

	if oldPosition == newPosition {
		return project, nil
	}

	project.Position = newPosition

	updatedProject, err := s.projectsRepository.ChangePosition(
		ctx, userID, project, oldPosition,
	)
	if err != nil {
		return domain.Project{}, fmt.Errorf(
			"error changing position in repository: %w", err,
		)
	}

	return updatedProject, nil
}
