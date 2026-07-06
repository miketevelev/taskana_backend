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

	CreateHeading(
		ctx context.Context,
		userID uuid.UUID,
		heading domain.Heading,
	) (domain.Heading, error)
}

func NewHeadingService(
	headingRepository HeadingRepository,
) *HeadingService {
	return &HeadingService{
		headingRepository: headingRepository,
	}
}
