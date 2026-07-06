package heading_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (r *HeadingRepository) GetHeadings(
	ctx context.Context,
	userID uuid.UUID,
	limit *int,
	offset *int,
) ([]domain.Heading, error) {
	query := `
		SELECT 
			h.id, h.version, h.user_id, h.project_id, h.title, 
			h.position, h.created_at, h.updated_at
		FROM taskana.headings h
		WHERE h.user_id = $1
		ORDER BY 
			h.project_id ASC, 
			h.position ASC, 
			h.created_at ASC 
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.pool.Query(
		ctx,
		query,
		userID,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("select headings: %w", err)
	}
	defer rows.Close()

	var headingModels []HeadingModel

	for rows.Next() {
		headingModel, err := scanHeading(rows)
		if err != nil {
			return nil, fmt.Errorf("scan headings from db: %w", err)
		}
		headingModels = append(headingModels, headingModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	headingDomains := headingDomainsFromModels(headingModels)

	return headingDomains, nil
}
