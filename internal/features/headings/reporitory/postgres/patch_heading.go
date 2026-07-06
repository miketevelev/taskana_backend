package heading_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
	core_postgres_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool"
)

func (r *HeadingRepository) PatchHeading(
	ctx context.Context,
	userID uuid.UUID,
	heading domain.Heading,
) (domain.Heading, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE taskana.headings
		SET 
			title = $1,
			updated_at = $2,
			version = version + 1
		WHERE id = $3 AND user_id = $4 AND version = $5
		RETURNING id, version, user_id, project_id, title, position, created_at, 
		       updated_at
	`

	oldVersion := heading.Version

	row := r.pool.QueryRow(
		ctx, query,
		heading.Title,
		time.Now().UTC(),
		heading.ID,
		userID,
		oldVersion,
	)

	headingModel, err := scanHeading(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Heading{}, fmt.Errorf(
				"heading with id='%s' concurrently accessed: %w",
				heading.ID,
				core_errors.ErrConflict,
			)
		}
		return domain.Heading{}, fmt.Errorf(
			"patch heading repository: %w", err,
		)
	}

	return headingDomainFromModel(headingModel), nil
}
