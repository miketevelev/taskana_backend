package projects_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *ProjectService) GetProject(
	ctx context.Context,
	userID uuid.UUID,
	projectID uuid.UUID,
) (domain.Project, error) {
	project, err := s.projectsRepository.GetProject(ctx, userID, projectID)
	if err != nil {
		return domain.Project{}, fmt.Errorf(
			"error getting project: %w", err,
		)
	}

	return project, nil
}
