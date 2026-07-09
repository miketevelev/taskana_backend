package headings_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *HeadingService) DeleteHeading(
	ctx context.Context,
	userID uuid.UUID,
	headingID uuid.UUID,
) error {
	if err := s.headingRepository.DeleteHeading(
		ctx, userID, headingID,
	); err != nil {
		return fmt.Errorf(
			"delete heading: %w", err,
		)
	}

	return nil
}
