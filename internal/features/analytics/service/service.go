package analytics_service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type AnalyticsFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
	ProjectID *uuid.UUID
}

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
