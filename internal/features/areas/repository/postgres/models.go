package areas_postgres_repository

import (
	"fmt"
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

func areaDomainsFromModels(areas []AreaModel) []domain.Area {
	if len(areas) == 0 {
		return []domain.Area{}
	}
	areaDomains := make([]domain.Area, len(areas))

	for i, area := range areas {
		areaDomains[i] = domain.Area{
			ID:        area.ID,
			Version:   area.Version,
			UserID:    area.UserID,
			Title:     area.Title,
			Position:  area.Position,
			CreatedAt: area.CreatedAt,
			UpdatedAt: area.UpdatedAt,
		}
	}

	return areaDomains
}

func scanArea(row interface{ Scan(dest ...any) error }) (AreaModel, error) {
	var areaModel AreaModel
	err := row.Scan(
		&areaModel.ID,
		&areaModel.Version,
		&areaModel.UserID,
		&areaModel.Title,
		&areaModel.Position,
		&areaModel.CreatedAt,
		&areaModel.UpdatedAt,
	)
	if err != nil {
		return AreaModel{}, fmt.Errorf("scan area: %w", err)
	}
	return areaModel, nil
}
