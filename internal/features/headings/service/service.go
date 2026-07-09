package headings_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type HeadingService struct {
	headingRepository HeadingRepository
}

type HeadingRepository interface {
	GetHeading(
		ctx context.Context,
		userID uuid.UUID,
		headingID uuid.UUID,
	) (domain.Heading, error)

	GetHeadings(
		ctx context.Context,
		userID uuid.UUID,
		limit *int,
		offset *int,
	) ([]domain.Heading, error)

	CreateHeading(
		ctx context.Context,
		userID uuid.UUID,
		heading domain.Heading,
	) (domain.Heading, error)

	ChangePosition(
		ctx context.Context,
		userID uuid.UUID,
		heading domain.Heading,
		oldPosition int,
	) (domain.Heading, error)

	PatchHeading(
		ctx context.Context,
		userID uuid.UUID,
		heading domain.Heading,
	) (domain.Heading, error)

	DeleteHeading(
		ctx context.Context,
		userID uuid.UUID,
		headingID uuid.UUID,
	) error
}

func NewHeadingService(
	headingRepository HeadingRepository,
) *HeadingService {
	return &HeadingService{
		headingRepository: headingRepository,
	}
}
