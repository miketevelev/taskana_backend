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

func areaDTOFromDomain(t domain.Area) AreaDTOResponse {
	return AreaDTOResponse{
		ID:        t.ID,
		Version:   t.Version,
		UserID:    t.UserID,
		Title:     t.Title,
		Position:  t.Position,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func areasDTOsFromDomains(areas []domain.Area) []AreaDTOResponse {
	result := make([]AreaDTOResponse, len(areas))
	for i, t := range areas {
		result[i] = areaDTOFromDomain(t)
	}
	return result
}
