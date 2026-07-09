package headings_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *HeadingService) GetHeading(
	ctx context.Context,
	userID uuid.UUID,
	headingID uuid.UUID,
) (domain.Heading, error) {
	heading, err := s.headingRepository.GetHeading(ctx, userID, headingID)
	if err != nil {
		return domain.Heading{}, fmt.Errorf(
			"error getting heading: %w", err,
		)
	}

	return heading, nil
}
