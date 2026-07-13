package areas_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain/area"
)

func (s *AreasService) GetArea(
	ctx context.Context,
	userID uuid.UUID,
	areaID uuid.UUID,
) (domain_area.Area, error) {
	area, err := s.areasRepository.GetArea(ctx, userID, areaID)
	if err != nil {
		return domain_area.Area{}, fmt.Errorf(
			"error getting area: %w", err,
		)
	}

	return area, nil
}
