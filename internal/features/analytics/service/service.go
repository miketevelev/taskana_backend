package analytics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type AnalyticsService struct {
	analyticsRepository AnalyticsRepository
}

type AnalyticsRepository interface {
	GetTimeByProject(
		ctx context.Context,
		userID uuid.UUID,
		startDate, endDate time.Time,
		projectID *uuid.UUID,
	) ([]domain.ProjectTimeSlice, error)

	GetFocusTimeline(
		ctx context.Context,
		userID uuid.UUID,
		startDate, endDate time.Time,
	) ([]domain.DailyFocusSlice, error)

	GetEstimationAccuracy(
		ctx context.Context,
		userID uuid.UUID,
		startDate, endDate time.Time,
		projectID *uuid.UUID,
	) ([]domain.EstimationAccuracyRow, error)
}

func NewAnalyticsService(
	analyticsRepository AnalyticsRepository,
) *AnalyticsService {
	return &AnalyticsService{
		analyticsRepository: analyticsRepository,
	}
}

func (s *AnalyticsService) getAnalyticsForRange(
	ctx context.Context,
	userID uuid.UUID,
	start, end time.Time,
	projectID *uuid.UUID,
) (domain.AnalyticsResponse, error) {
	timeByProject, err := s.analyticsRepository.GetTimeByProject(
		ctx, userID, start, end, projectID,
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
		ctx, userID, start, end, projectID,
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

func dayBounds(day time.Time) (start, end time.Time) {
	start = time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	end = start.Add(24*time.Hour - time.Nanosecond)
	return start, end
}

func weekBounds(reference time.Time, capToNow bool) (start, end time.Time) {
	weekday := int(reference.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday -> 7, so Monday is always offset 1
	}
	monday := reference.AddDate(0, 0, -(weekday - 1))
	start = time.Date(
		monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, time.UTC,
	)
	end = start.AddDate(0, 0, 7).Add(-time.Nanosecond)

	if capToNow {
		now := time.Now().UTC()
		if end.After(now) {
			end = now
		}
	}
	return start, end
}

func monthBounds(reference time.Time, capToNow bool) (start, end time.Time) {
	start = time.Date(
		reference.Year(), reference.Month(), 1, 0, 0, 0, 0, time.UTC,
	)
	end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	if capToNow {
		now := time.Now().UTC()
		if end.After(now) {
			end = now
		}
	}
	return start, end
}

func yearBounds(year int, capToNow bool) (start, end time.Time) {
	start = time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	end = start.AddDate(1, 0, 0).Add(-time.Nanosecond)

	if capToNow {
		now := time.Now().UTC()
		if end.After(now) {
			end = now
		}
	}
	return start, end
}

type AnalyticsFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
	ProjectID *uuid.UUID
}
