package areas_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type AreaService struct {
	areaRepository AreaRepository
}

type AreaRepository interface {
	CreateArea(
		ctx context.Context,
		userID uuid.UUID,
		area domain.Area,
	) (domain.Area, error)
}

func NewAreaService(
	areaRepository AreaRepository,
) *AreaService {
	return &AreaService{
		areaRepository: areaRepository,
	}
}
