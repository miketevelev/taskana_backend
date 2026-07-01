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
	GetArea(
		ctx context.Context,
		userID uuid.UUID,
		areaID uuid.UUID,
	) (domain.Area, error)

	GetAreas(
		ctx context.Context,
		userID uuid.UUID,
		limit *int,
		offset *int,
	) ([]domain.Area, error)

	CreateArea(
		ctx context.Context,
		userID uuid.UUID,
		area domain.Area,
	) (domain.Area, error)

	DeleteArea(
		ctx context.Context,
		userID uuid.UUID,
		areaID uuid.UUID,
	) error
}

func NewAreaService(
	areaRepository AreaRepository,
) *AreaService {
	return &AreaService{
		areaRepository: areaRepository,
	}
}
