package areas_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *AreasService) DeleteArea(
	ctx context.Context,
	userID uuid.UUID,
	areaID uuid.UUID,
) error {
	if err := s.areasRepository.DeleteArea(ctx, userID, areaID); err != nil {
		return fmt.Errorf("delete area: %w", err)
	}

	return nil
}
