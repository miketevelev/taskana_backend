package areas_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	domain_area "github.com/miketevelev/taskana_backend/internal/core/domain/area"
)

func (s *AreasService) CreateArea(
	ctx context.Context,
	userID uuid.UUID,
	area domain_area.Area,
) (domain_area.Area, error) {
	if area.ID == uuid.Nil {
		area.ID = uuid.New()
	}

	if area.Version == -1 {
		area.Version = 1
	}

	if err := area.Validate(); err != nil {
		return domain_area.Area{},
			fmt.Errorf("area validation failed: %w", err)
	}

	createdArea, err := s.areasRepository.CreateArea(ctx, userID, area)
	if err != nil {
		return domain_area.Area{}, fmt.Errorf("create area failed: %w", err)
	}

	return createdArea, nil
}
