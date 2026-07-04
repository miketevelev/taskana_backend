package projects_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *ProjectService) DeleteProject(
	ctx context.Context,
	userID uuid.UUID,
	projectID uuid.UUID,
) error {
	if err := s.projectsRepository.DeleteProject(
		ctx, userID, projectID,
	); err != nil {
		return fmt.Errorf(
			"delete project: %w", err,
		)
	}

	return nil
}
