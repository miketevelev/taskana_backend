package headings_transport_http

import (
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type HeadingDTOResponse struct {
	ID        uuid.UUID `json:"id"`
	Version   int       `json:"version"`
	UserID    uuid.UUID `json:"user_id"`
	ProjectID uuid.UUID `json:"project_id"`
	Title     string    `json:"title"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func headingDTOFromDomain(heading domain.Heading) HeadingDTOResponse {
	return HeadingDTOResponse{
		ID:        heading.ID,
		Version:   heading.Version,
		UserID:    heading.UserID,
		ProjectID: heading.ProjectID,
		Title:     heading.Title,
		Position:  heading.Position,
		CreatedAt: heading.CreatedAt,
		UpdatedAt: heading.UpdatedAt,
	}
}

func headingDTOsFromDomains(headings []domain.Heading) []HeadingDTOResponse {
	result := make([]HeadingDTOResponse, len(headings))
	for i, heading := range headings {
		result[i] = headingDTOFromDomain(heading)
	}
	return result
}
