package headings_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *HeadingService) PatchHeading(
	ctx context.Context,
	userID uuid.UUID,
	headingID uuid.UUID,
	patch domain.HeadingPatch,
) (domain.Heading, error) {
	heading, err := s.headingRepository.GetHeading(ctx, userID, headingID)
	if err != nil {
		return domain.Heading{}, fmt.Errorf(
			"error while fetching heading: %w", err,
		)
	}

	originalTitle := heading.Title

	if err := heading.ApplyPatch(patch); err != nil {
		return domain.Heading{}, fmt.Errorf(
			"error applying patch to heading: %w", err,
		)
	}

	if heading.Title == originalTitle {
		return heading, nil
	}

	patchedHeading, err := s.headingRepository.PatchHeading(
		ctx, userID, heading,
	)
	if err != nil {
		return domain.Heading{}, fmt.Errorf(
			"error while saving patched heading: %w", err,
		)
	}

	return patchedHeading, nil
}
