package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
)

type PomodoroSession struct {
	ID              uuid.UUID `json:"id"`
	Version         int       `json:"version"`
	UserID          uuid.UUID `json:"user_id"`
	TaskID          uuid.UUID `json:"task_id"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	DurationSeconds int       `json:"duration_seconds"`
	IsInterrupted   bool      `json:"is_interrupted"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func NewPomodoroSession(
	id uuid.UUID,
	version int,
	userID uuid.UUID,
	taskID uuid.UUID,
	startTime time.Time,
	endTime time.Time,
	durationSeconds int,
	isInterrupted bool,
	createdAt time.Time,
	updatedAt time.Time,
) PomodoroSession {
	return PomodoroSession{
		ID:              id,
		Version:         version,
		UserID:          userID,
		TaskID:          taskID,
		StartTime:       startTime,
		EndTime:         endTime,
		DurationSeconds: durationSeconds,
		IsInterrupted:   isInterrupted,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
}

func NewPomodoroSessionUninitialized(
	id uuid.UUID,
	userID uuid.UUID,
	taskID uuid.UUID,
	startTime time.Time,
	endTime time.Time,
	durationSeconds int,
	isInterrupted bool,
) PomodoroSession {
	now := time.Now().UTC()
	return NewPomodoroSession(
		id,
		UninitializedVersion,
		userID,
		taskID,
		startTime,
		endTime,
		durationSeconds,
		isInterrupted,
		now,
		now,
	)
}

func (p *PomodoroSession) Validate() error {
	if p.UserID == uuid.Nil {
		return fmt.Errorf(
			"user_id cannot be an empty UUID: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.TaskID == uuid.Nil {
		return fmt.Errorf(
			"task_id cannot be an empty UUID: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.StartTime.IsZero() {
		return fmt.Errorf(
			"start_time cannot be zero: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.EndTime.IsZero() {
		return fmt.Errorf(
			"end_time cannot be zero: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.EndTime.Before(p.StartTime) {
		return fmt.Errorf(
			"end_time cannot be before start_time: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.DurationSeconds < 0 {
		return fmt.Errorf(
			"duration_seconds cannot be negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.CreatedAt.IsZero() || p.UpdatedAt.IsZero() {
		return fmt.Errorf(
			"timestamps cannot be zero: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.UpdatedAt.Before(p.CreatedAt) {
		return fmt.Errorf(
			"updated_at cannot be before created_at: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.Version < 0 {
		p.Version = 0
	}

	return nil
}

func (p *PomodoroSession) ApplyPatch(patch PomodoroSessionPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate pomodoro session patch: %w", err)
	}

	tmp := *p

	if patch.StartTime.Set {
		tmp.StartTime = *patch.StartTime.Value
	}

	if patch.EndTime.Set {
		tmp.EndTime = *patch.EndTime.Value
	}

	if patch.DurationSeconds.Set {
		tmp.DurationSeconds = *patch.DurationSeconds.Value
	}

	if patch.IsInterrupted.Set {
		tmp.IsInterrupted = *patch.IsInterrupted.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate pomodoro session after patch: %w", err)
	}

	*p = tmp

	return nil
}

type PomodoroSessionPatch struct {
	StartTime       Nullable[time.Time]
	EndTime         Nullable[time.Time]
	DurationSeconds Nullable[int]
	IsInterrupted   Nullable[bool]
}

func NewPomodoroSessionPatch(
	startTime Nullable[time.Time],
	endTime Nullable[time.Time],
	durationSeconds Nullable[int],
	isInterrupted Nullable[bool],
) PomodoroSessionPatch {
	return PomodoroSessionPatch{
		StartTime:       startTime,
		EndTime:         endTime,
		DurationSeconds: durationSeconds,
		IsInterrupted:   isInterrupted,
	}
}

func (p *PomodoroSessionPatch) Validate() error {
	if p.StartTime.Set && p.StartTime.Value == nil {
		return fmt.Errorf(
			"'StartTime' can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	if p.StartTime.Set && p.StartTime.Value != nil {
		if (*p.StartTime.Value).IsZero() {
			return fmt.Errorf(
				"patch 'StartTime' cannot be a zero time: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if p.EndTime.Set && p.EndTime.Value == nil {
		return fmt.Errorf(
			"'EndTime' can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	if p.EndTime.Set && p.EndTime.Value != nil {
		if (*p.EndTime.Value).IsZero() {
			return fmt.Errorf(
				"patch 'EndTime' cannot be a zero time: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if p.DurationSeconds.Set && p.DurationSeconds.Value == nil {
		return fmt.Errorf(
			"'DurationSeconds' can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	if p.DurationSeconds.Set && p.DurationSeconds.Value != nil {
		if *p.DurationSeconds.Value < 0 {
			return fmt.Errorf(
				"patch 'DurationSeconds' cannot be negative: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if p.IsInterrupted.Set && p.IsInterrupted.Value == nil {
		return fmt.Errorf(
			"'IsInterrupted' can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}
