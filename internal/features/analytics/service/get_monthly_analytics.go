package analytics_service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type MonthlyAnalyticsFilter struct {
	// Month may be any date that falls within the desired month.
	// If nil, the current month is used (capped to "now", see monthBounds).
	Month     *time.Time
	ProjectID *uuid.UUID
}

func (s *AnalyticsService) GetMonthlyAnalytics(
	ctx context.Context,
	userID uuid.UUID,
	filter MonthlyAnalyticsFilter,
) (domain.AnalyticsResponse, error) {
	reference := time.Now().UTC()
	capToNow := filter.Month == nil
	if filter.Month != nil {
		reference = *filter.Month
	}

	start, end := monthBounds(reference, capToNow)

	return s.getAnalyticsForRange(ctx, userID, start, end, filter.ProjectID)
}
