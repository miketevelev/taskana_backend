package areas_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *AreaService) DeleteArea(
	ctx context.Context,
	userID uuid.UUID,
	areaID uuid.UUID,
) error {
	if err := s.areaRepository.DeleteArea(ctx, userID, areaID); err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	return nil
}
