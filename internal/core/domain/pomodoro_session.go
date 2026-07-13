package domain

import (
	"time"

	"github.com/google/uuid"
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
