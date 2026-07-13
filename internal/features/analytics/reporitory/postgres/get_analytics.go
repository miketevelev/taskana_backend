package analytics_postgres_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (r *AnalyticsRepository) GetTimeByProject(
	ctx context.Context,
	userID uuid.UUID,
	startDate, endDate time.Time,
	projectID *uuid.UUID,
) ([]domain.ProjectTimeSlice, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT
   			COALESCE(p.id, '00000000-0000-0000-0000-000000000000'::uuid) AS project_id, 
   			COALESCE(p.title, 'No Project') AS project_title,
   			SUM(ps.duration_seconds) AS total_seconds
       	FROM taskana.pomodoro_sessions ps
       	JOIN taskana.tasks t ON ps.task_id = t.id
       	LEFT JOIN taskana.projects p ON t.project_id = p.id -- <-- LEFT JOIN
       	WHERE ps.user_id = $1
         	AND ps.start_time BETWEEN $2 AND $3
         	AND t.is_time_tracked = true
	`
	args := []any{userID, startDate, endDate}

	if projectID != nil {
		query += ` AND p.id = $4`
		args = append(args, *projectID)
	}

	query += ` GROUP BY p.id, p.title ORDER BY total_seconds DESC`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("time by project: %w", err)
	}
	defer rows.Close()

	var result []domain.ProjectTimeSlice
	for rows.Next() {
		var slice domain.ProjectTimeSlice
		if err := rows.Scan(
			&slice.ProjectID, &slice.ProjectTitle, &slice.TotalSeconds,
		); err != nil {
			return nil, fmt.Errorf("scan project time: %w", err)
		}
		result = append(result, slice)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate project time: %w", err)
	}

	return result, nil
}

func (r *AnalyticsRepository) GetFocusTimeline(
	ctx context.Context,
	userID uuid.UUID,
	startDate, endDate time.Time,
) ([]domain.DailyFocusSlice, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT
			DATE_TRUNC('day', ps.start_time) AS focus_day,
			SUM(ps.duration_seconds) AS daily_seconds,
			COUNT(ps.id) AS sessions_count
		FROM taskana.pomodoro_sessions ps
		WHERE ps.user_id = $1
		  AND ps.start_time BETWEEN $2 AND $3
		GROUP BY focus_day
		ORDER BY focus_day ASC
	`

	rows, err := r.pool.Query(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("focus timeline: %w", err)
	}
	defer rows.Close()

	var result []domain.DailyFocusSlice
	for rows.Next() {
		var slice domain.DailyFocusSlice
		if err := rows.Scan(
			&slice.FocusDay, &slice.DailySeconds, &slice.SessionsCount,
		); err != nil {
			return nil, fmt.Errorf("scan focus timeline: %w", err)
		}
		result = append(result, slice)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate focus timeline: %w", err)
	}

	return result, nil
}

func (r *AnalyticsRepository) GetEstimationAccuracy(
	ctx context.Context,
	userID uuid.UUID,
	startDate, endDate time.Time,
	projectID *uuid.UUID,
) ([]domain.EstimationAccuracyRow, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
       SELECT
          t.id AS task_id,
          t.title AS task_title,
          t.estimated_pomodoros AS planned_pomodoros,
          COUNT(ps.id) AS actual_pomodoros,
          CASE
             WHEN t.estimated_pomodoros = 0 THEN 0
             ELSE ROUND((COUNT(ps.id)::numeric / t.estimated_pomodoros::numeric) * 100, 2)
          END AS accuracy_percentage
       FROM taskana.tasks t
       LEFT JOIN taskana.pomodoro_sessions ps
          ON t.id = ps.task_id AND ps.is_interrupted = false
       WHERE t.user_id = $1
         AND t.status = 'completed'
         AND t.is_time_tracked = true
         AND t.completed_at BETWEEN $2 AND $3 -- <-- Фильтруем по дате завершения задачи!
    `
	args := []any{userID, startDate, endDate}

	if projectID != nil {
		query += ` AND t.project_id = $4`
		args = append(args, *projectID)
	}

	query += ` GROUP BY t.id, t.title, t.estimated_pomodoros ORDER BY t.completed_at DESC`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("estimation accuracy: %w", err)
	}
	defer rows.Close()

	var result []domain.EstimationAccuracyRow
	for rows.Next() {
		var row domain.EstimationAccuracyRow
		if err := rows.Scan(
			&row.TaskID,
			&row.TaskTitle,
			&row.PlannedPomodoros,
			&row.ActualPomodoros,
			&row.AccuracyPercentage,
		); err != nil {
			return nil, fmt.Errorf("scan estimation accuracy: %w", err)
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate estimation accuracy: %w", err)
	}

	return result, nil
}
