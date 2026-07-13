package areas_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	domain_area "github.com/miketevelev/taskana_backend/internal/core/domain/area"
)

func (s *AreasService) PatchArea(
	ctx context.Context,
	userID uuid.UUID,
	areaID uuid.UUID,
	patch domain_area.AreaPatch,
) (domain_area.Area, error) {
	area, err := s.areasRepository.GetArea(ctx, userID, areaID)
	if err != nil {
		return domain_area.Area{}, fmt.Errorf(
			"error while fetching area: %w", err,
		)
	}

	if err := area.ApplyPatch(patch); err != nil {
		return domain_area.Area{}, fmt.Errorf(
			"error while applying patch to area: %w", err,
		)
	}

	patchedArea, err := s.areasRepository.PatchArea(ctx, userID, area)
	if err != nil {
		return domain_area.Area{}, fmt.Errorf(
			"error while saving patched area: %w", err,
		)
	}

	return patchedArea, nil

	return patchedArea, nil
}
