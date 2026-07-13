package timetracking_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *TimeTrackingService) Sync(
	ctx context.Context,
	userID uuid.UUID,
	sessions []domain.PomodoroSession,
) error {
	for i := range sessions {
		if sessions[i].Version == domain.UninitializedVersion {
			sessions[i].Version = 1
		}

		if err := sessions[i].Validate(); err != nil {
			return fmt.Errorf(
				"pomodoro session validation failed: %w", err,
			)
		}
	}

	if err := s.timeTrackingRepository.BulkInsertSessions(
		ctx, userID, sessions,
	); err != nil {
		return fmt.Errorf("bulk insert sessions: %w", err)
	}

	//if err := s.timeTrackingRepository.InvalidateAnalyticsCache(
	//	ctx, userID,
	//); err != nil {
	//	return fmt.Errorf("invalidate analytics cache: %w", err)
	//}

	return nil
}
