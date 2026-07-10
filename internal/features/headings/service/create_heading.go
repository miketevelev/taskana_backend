package headings_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *HeadingService) CreateHeading(
	ctx context.Context,
	userID uuid.UUID,
	heading domain.Heading,
) (domain.Heading, error) {
	if err := heading.Validate(); err != nil {
		return domain.Heading{}, fmt.Errorf(
			"heading validation failed: %w",
		)
	}

	createdHeading, err := s.headingRepository.CreateHeading(
		ctx,
		userID,
		heading,
	)
	if err != nil {
		return domain.Heading{}, fmt.Errorf(
			"failed to create heading: %w",
			err,
		)
	}

	return createdHeading, nil
}
