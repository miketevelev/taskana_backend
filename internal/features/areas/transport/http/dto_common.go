package areas_transport_http

import (
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type AreaDTOResponse struct {
	ID        uuid.UUID `json:"id"`
	Version   int       `json:"version"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func areaDTOFromDomain(area domain.Area) AreaDTOResponse {
	return AreaDTOResponse{
		ID:        area.ID,
		Version:   area.Version,
		UserID:    area.UserID,
		Title:     area.Title,
		Position:  area.Position,
		CreatedAt: area.CreatedAt,
		UpdatedAt: area.UpdatedAt,
	}
}

func areasDTOsFromDomains(areas []domain.Area) []AreaDTOResponse {
	result := make([]AreaDTOResponse, len(areas))
	for i, t := range areas {
		result[i] = areaDTOFromDomain(t)
	}
	return result
}
