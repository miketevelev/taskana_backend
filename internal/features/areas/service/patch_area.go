package areas_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *AreasService) PatchArea(
	ctx context.Context,
	userID uuid.UUID,
	areaID uuid.UUID,
	patch domain.AreaPatch,
) (domain.Area, error) {
	area, err := s.areasRepository.GetArea(ctx, userID, areaID)
	if err != nil {
		return domain.Area{}, fmt.Errorf("error while fetching area: %w", err)
	}

	if err := area.ApplyPatch(patch); err != nil {
		return domain.Area{}, fmt.Errorf(
			"error while applying patch to area: %w", err,
		)
	}

	patchedArea, err := s.areasRepository.PatchArea(ctx, userID, area)
	if err != nil {
		return domain.Area{}, fmt.Errorf(
			"error while saving patched area: %w", err,
		)
	}

	return patchedArea, nil

	return patchedArea, nil
}
