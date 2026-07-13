package timetracking_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type TimeTrackingService struct {
	timeTrackingRepository TimeTrackingRepository
}

type TimeTrackingRepository interface {
	BulkInsertSessions(
		ctx context.Context,
		userID uuid.UUID,
		sessions []domain.PomodoroSession,
	) error

	//InvalidateAnalyticsCache(
	//	ctx context.Context,
	//	userID uuid.UUID,
	//) error
}

func NewTimeTrackingService(
	timeTrackingRepository TimeTrackingRepository,
) *TimeTrackingService {
	return &TimeTrackingService{
		timeTrackingRepository: timeTrackingRepository,
	}
}
