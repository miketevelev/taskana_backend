package analytics_service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type YearlyAnalyticsFilter struct {
	// Year is the calendar year to report on. If nil, the current year is used
	// (capped to "now", see yearBounds).
	Year      *int
	ProjectID *uuid.UUID
}

func (s *AnalyticsService) GetYearlyAnalytics(
	ctx context.Context,
	userID uuid.UUID,
	filter YearlyAnalyticsFilter,
) (domain.AnalyticsResponse, error) {
	year := time.Now().UTC().Year()
	capToNow := filter.Year == nil
	if filter.Year != nil {
		year = *filter.Year
	}

	start, end := yearBounds(year, capToNow)

	return s.getAnalyticsForRange(ctx, userID, start, end, filter.ProjectID)
}
