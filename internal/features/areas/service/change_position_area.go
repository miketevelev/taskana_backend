package areas_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *AreasService) ChangePosition(
	ctx context.Context,
	userID uuid.UUID,
	areaID uuid.UUID,
	newPosition int,
) (domain.Area, error) {
	area, err := s.areasRepository.GetArea(ctx, userID, areaID)
	if err != nil {
		return domain.Area{}, fmt.Errorf("error getting area: %w", err)
	}

	oldPosition := area.Position

	if oldPosition == newPosition {
		return area, nil
	}

	area.Position = newPosition

	updatedArea, err := s.areasRepository.ChangePosition(
		ctx, userID, area, oldPosition,
	)
	if err != nil {
		return domain.Area{}, fmt.Errorf(
			"error changing position in repository: %w", err,
		)
	}

	return updatedArea, nil
}
