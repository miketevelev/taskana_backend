package areas_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *AreaService) CreateArea(
	ctx context.Context,
	userID uuid.UUID,
	area domain.Area,
) (domain.Area, error) {
	if area.ID == uuid.Nil {
		area.ID = uuid.New()
	}

	if area.Version == -1 {
		area.Version = 1
	}

	if err := area.Validate(); err != nil {
		return domain.Area{},
			fmt.Errorf("area validation failed: %w", err)
	}

	createdArea, err := s.areaRepository.CreateArea(ctx, userID, area)
	if err != nil {
		return domain.Area{}, fmt.Errorf("create area failed: %w", err)
	}

	return createdArea, nil
}
