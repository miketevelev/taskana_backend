package projects_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain/project"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
)

func (s *ProjectService) GetProjects(
	ctx context.Context,
	userID uuid.UUID,
	limit *int,
	offset *int,
) ([]domain_project.Project, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf(
			"limit must be non-negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf(
			"offset must be non-negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	projects, err := s.projectsRepository.GetProjects(
		ctx, userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("get projects from repository: %w", err)
	}

	return projects, nil
}
