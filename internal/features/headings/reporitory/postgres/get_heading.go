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

func (r *HeadingRepository) GetHeading(
	ctx context.Context,
	userID uuid.UUID,
	headingID uuid.UUID,
) (domain.Heading, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, user_id, project_id, title, position, created_at, 
		       updated_at
		FROM taskana.headings
		WHERE id = $1 AND user_id = $2
	`

	row := r.pool.QueryRow(ctx, query, headingID, userID)

	headingModel, err := scanHeading(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Heading{}, fmt.Errorf(
				"heading with id %s not found: %w",
				headingID,
				core_errors.ErrNotFound,
			)
		}

		return domain.Heading{}, fmt.Errorf("scan error %w", err)
	}

	headingDomain := headingDomainFromModel(headingModel)

	return headingDomain, nil
}
