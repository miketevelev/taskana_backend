package areas_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain/area"
)

type AreasService struct {
	areasRepository AreasRepository
}

type AreasRepository interface {
	GetArea(
		ctx context.Context,
		userID uuid.UUID,
		areaID uuid.UUID,
	) (domain_area.Area, error)

	GetAreas(
		ctx context.Context,
		userID uuid.UUID,
		limit *int,
		offset *int,
	) ([]domain_area.Area, error)

	CreateArea(
		ctx context.Context,
		userID uuid.UUID,
		area domain_area.Area,
	) (domain_area.Area, error)

	ChangePosition(
		ctx context.Context,
		userID uuid.UUID,
		area domain_area.Area,
		oldPosition int,
	) (domain_area.Area, error)

	PatchArea(
		ctx context.Context,
		userID uuid.UUID,
		area domain_area.Area,
	) (domain_area.Area, error)

	PatchAreaWithReordering(
		ctx context.Context,
		userID uuid.UUID,
		area domain_area.Area,
		oldPos int,
	) (domain_area.Area, error)

	DeleteArea(
		ctx context.Context,
		userID uuid.UUID,
		areaID uuid.UUID,
	) error
}

func NewAreasService(
	areaRepository AreasRepository,
) *AreasService {
	return &AreasService{
		areasRepository: areaRepository,
	}
}
