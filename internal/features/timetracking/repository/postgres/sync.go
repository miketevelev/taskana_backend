package timetracking_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
)

func (r *TimeTrackingRepository) BulkInsertSessions(
	ctx context.Context,
	userID uuid.UUID,
	sessions []domain.PomodoroSession,
) error {
	if len(sessions) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	uniqueTasks := make(map[uuid.UUID]struct{})
	for _, s := range sessions {
		if s.UserID != userID {
			return fmt.Errorf("session user mismatch")
		}
		uniqueTasks[s.TaskID] = struct{}{}
	}

	taskIDs := make([]uuid.UUID, 0, len(uniqueTasks))
	for id := range uniqueTasks {
		taskIDs = append(taskIDs, id)
	}

	checkQuery := `
		SELECT COUNT(id) FROM taskana.tasks 
		WHERE id = ANY($1) AND user_id = $2
	`
	var foundCount int
	err := r.pool.QueryRow(ctx, checkQuery, taskIDs, userID).Scan(&foundCount)
	if err != nil {
		return fmt.Errorf("verify tasks ownership: %w", err)
	}
	if foundCount != len(taskIDs) {
		return fmt.Errorf(
			"some tasks not found or do not belong to user: %w",
			core_errors.ErrNotFound,
		)
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	insertQuery := `
		INSERT INTO taskana.pomodoro_sessions (
			id, version, user_id, task_id, start_time, end_time, 
			duration_seconds, is_interrupted, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO NOTHING
	`

	for _, session := range sessions {
		_, err = tx.Exec(
			ctx, insertQuery,
			session.ID, session.Version, userID, session.TaskID,
			session.StartTime, session.EndTime, session.DurationSeconds,
			session.IsInterrupted, session.CreatedAt, session.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("insert session %s: %w", session.ID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
