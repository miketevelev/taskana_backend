package headings_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *HeadingService) ChangePosition(
	ctx context.Context,
	userID uuid.UUID,
	headingID uuid.UUID,
	newPosition int,
) (domain.Heading, error) {
	heading, err := s.headingRepository.GetHeading(ctx, userID, headingID)
	if err != nil {
		return domain.Heading{}, fmt.Errorf(
			"error getting project: %w", err,
		)
	}

	oldPosition := heading.Position

	if oldPosition == newPosition {
		return heading, nil
	}

	heading.Position = newPosition

	updatedHeading, err := s.headingRepository.ChangePosition(
		ctx, userID, heading, oldPosition,
	)
	if err != nil {
		return domain.Heading{}, fmt.Errorf(
			"error changing position in repository: %w", err,
		)
	}

	return updatedHeading, nil
}
