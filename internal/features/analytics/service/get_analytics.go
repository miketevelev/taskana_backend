package analytics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *AnalyticsService) GetAnalytics(
	ctx context.Context,
	userID uuid.UUID,
	filter AnalyticsFilter,
) (domain.AnalyticsResponse, error) {
	start := time.Now().UTC().AddDate(0, -1, 0)
	end := time.Now().UTC()
	if filter.StartDate != nil {
		start = *filter.StartDate
	}
	if filter.EndDate != nil {
		end = filter.EndDate.Add(24*time.Hour - time.Nanosecond)
	}

	timeByProject, err := s.analyticsRepository.GetTimeByProject(
		ctx, userID, start, end, filter.ProjectID,
	)
	if err != nil {
		return domain.AnalyticsResponse{}, fmt.Errorf(
			"time by project: %w", err,
		)
	}

	focusTimeline, err := s.analyticsRepository.GetFocusTimeline(
		ctx, userID, start, end,
	)
	if err != nil {
		return domain.AnalyticsResponse{}, fmt.Errorf("focus timeline: %w", err)
	}

	estimation, err := s.analyticsRepository.GetEstimationAccuracy(
		ctx, userID, start, end, filter.ProjectID,
	)
	if err != nil {
		return domain.AnalyticsResponse{}, fmt.Errorf(
			"estimation accuracy: %w", err,
		)
	}

	return domain.AnalyticsResponse{
		TimeByProject:      timeByProject,
		FocusTimeline:      focusTimeline,
		EstimationAccuracy: estimation,
	}, nil
}
