package analytics_service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type WeeklyAnalyticsFilter struct {
	// Week may be any date that falls within the desired week (Mon..Sun).
	// If nil, the current week is used (capped to "now", see weekBounds).
	Week      *time.Time
	ProjectID *uuid.UUID
}

func (s *AnalyticsService) GetWeeklyAnalytics(
	ctx context.Context,
	userID uuid.UUID,
	filter WeeklyAnalyticsFilter,
) (domain.AnalyticsResponse, error) {
	reference := time.Now().UTC()
	capToNow := filter.Week == nil
	if filter.Week != nil {
		reference = *filter.Week
	}

	start, end := weekBounds(reference, capToNow)

	return s.getAnalyticsForRange(ctx, userID, start, end, filter.ProjectID)
}
