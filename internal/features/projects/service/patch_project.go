package projects_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *ProjectService) PatchProject(
	ctx context.Context,
	userID uuid.UUID,
	projectID uuid.UUID,
	patch domain.ProjectPatch,
) (domain.Project, error) {
	project, err := s.projectsRepository.GetProject(ctx, userID, projectID)
	if err != nil {
		return domain.Project{}, fmt.Errorf(
			"error while fetching project: %w", err,
		)
	}

	oldAreaID := project.AreaID
	oldPosition := project.Position

	if err := project.ApplyPatch(patch); err != nil {
		return domain.Project{}, fmt.Errorf(
			"error while applying patch to project: %w", err,
		)
	}

	areaChanged := false
	if patch.AreaID.Set {
		if (oldAreaID == nil && project.AreaID != nil) ||
			(oldAreaID != nil && project.AreaID == nil) ||
			(oldAreaID != nil && project.AreaID != nil && *oldAreaID != *project.AreaID) {
			areaChanged = true
		}
	}

	var patchedProject domain.Project

	if areaChanged {
		patchedProject, err = s.projectsRepository.PatchProjectWithAreaChange(
			ctx, userID, project, oldPosition, oldAreaID,
		)
	} else {
		patchedProject, err = s.projectsRepository.PatchProject(
			ctx, userID, project,
		)
	}

	if err != nil {
		return domain.Project{}, fmt.Errorf(
			"error while saving patched project: %w", err,
		)
	}

	return patchedProject, nil
}
