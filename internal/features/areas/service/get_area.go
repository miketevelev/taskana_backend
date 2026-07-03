package areas_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *AreaService) GetArea(
	ctx context.Context,
	userID uuid.UUID,
	areaID uuid.UUID,
) (domain.Area, error) {
	area, err := s.areaRepository.GetArea(ctx, userID, areaID)
	if err != nil {
		return domain.Area{}, fmt.Errorf(
			"failed to decode area request: %w", err,
		)
	}

	return area, nil
}
