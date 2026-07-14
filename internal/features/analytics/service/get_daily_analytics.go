package analytics_service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type DailyAnalyticsFilter struct {
	Day       *time.Time
	ProjectID *uuid.UUID
}

func (s *AnalyticsService) GetDailyAnalytics(
	ctx context.Context,
	userID uuid.UUID,
	filter DailyAnalyticsFilter,
) (domain.AnalyticsResponse, error) {
	day := time.Now().UTC()
	if filter.Day != nil {
		day = *filter.Day
	}

	start, end := dayBounds(day)

	return s.getAnalyticsForRange(ctx, userID, start, end, filter.ProjectID)
}
