package heading_postgres_repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type HeadingModel struct {
	ID        uuid.UUID `json:"id"`
	Version   int       `json:"version"`
	UserID    uuid.UUID `json:"user_id"`
	ProjectID uuid.UUID `json:"project_id"`
	Title     string    `json:"title"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func headingDomainFromModel(headingModel HeadingModel) domain.Heading {
	return domain.Heading{
		ID:        headingModel.ID,
		Version:   headingModel.Version,
		UserID:    headingModel.UserID,
		ProjectID: headingModel.ProjectID,
		Title:     headingModel.Title,
		Position:  headingModel.Position,
		CreatedAt: headingModel.CreatedAt,
		UpdatedAt: headingModel.UpdatedAt,
	}
}

func headingDomainsFromModels(headings []HeadingModel) []domain.Heading {
	if len(headings) == 0 {
		return []domain.Heading{}
	}
	headingDomains := make([]domain.Heading, len(headings))

	for i, heading := range headings {
		headingDomains[i] = domain.Heading{
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

	return headingDomains
}

func scanHeading(row interface{ Scan(dest ...any) error }) (
	HeadingModel,
	error,
) {
	var headingModel HeadingModel
	err := row.Scan(
		&headingModel.ID,
		&headingModel.Version,
		&headingModel.UserID,
		&headingModel.ProjectID,
		&headingModel.Title,
		&headingModel.Position,
		&headingModel.CreatedAt,
		&headingModel.UpdatedAt,
	)
	if err != nil {
		return HeadingModel{}, fmt.Errorf("scan heading: %w", err)
	}
	return headingModel, nil
}
