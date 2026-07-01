package areas_postgres_repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type AreaModel struct {
	ID        uuid.UUID `json:"id"`
	Version   int       `json:"version"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func areaFromModel(areaModel AreaModel) domain.Area {
	return domain.Area{
		ID:        areaModel.ID,
		Version:   areaModel.Version,
		UserID:    areaModel.UserID,
		Title:     areaModel.Title,
		Position:  areaModel.Position,
		CreatedAt: areaModel.CreatedAt,
		UpdatedAt: areaModel.UpdatedAt,
	}
}

func areaDomainFromModel(areaModel AreaModel) domain.Area {
	return domain.NewArea(
		areaModel.ID,
		areaModel.Version,
		areaModel.UserID,
		areaModel.Title,
		areaModel.Position,
		areaModel.CreatedAt,
		areaModel.UpdatedAt,
	)
}
