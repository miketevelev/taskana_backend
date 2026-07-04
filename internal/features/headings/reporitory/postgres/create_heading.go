package heading_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
	core_postgres_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool"
)

func (r *HeadingRepository) CreateHeading(
	ctx context.Context,
	userID uuid.UUID,
	heading domain.Heading,
) (domain.Heading, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO taskana.headings (id, version, user_id, project_id, title, 
		                              position, created_at, updated_at)
		VALUES (
			$1, $2, $3, $4, $5,
			(
				SELECT COALESCE(MAX(position), 0) + 1 
				FROM taskana.headings 
				WHERE user_id = $3 AND project_id = $4
			),
			$6, $7
		)
		RETURNING id, version, user_id, project_id, title, position, 
		    created_at, updated_at
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		heading.ID,
		heading.Version,
		userID,
		heading.ProjectID,
		heading.Title,
		heading.CreatedAt,
		heading.UpdatedAt,
	)

	headingModel, err := scanHeading(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolateForeignKey) {
			return domain.Heading{}, fmt.Errorf(
				"user or project not found for new heading: %w",
				core_errors.ErrNotFound,
			)
		}

		return domain.Heading{}, fmt.Errorf(
			"scan heading from db: %w", err,
		)
	}

	headingDomain := headingDomainFromModel(headingModel)

	return headingDomain, nil
}
