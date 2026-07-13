package timetracking_transport_http

import (
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type SessionDTOResponse struct {
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

func sessionDTOFromDomain(session domain.PomodoroSession) SessionDTOResponse {
	return SessionDTOResponse{
		ID:              session.ID,
		Version:         session.Version,
		UserID:          session.UserID,
		TaskID:          session.TaskID,
		StartTime:       session.StartTime,
		EndTime:         session.EndTime,
		DurationSeconds: session.DurationSeconds,
		IsInterrupted:   session.IsInterrupted,
		CreatedAt:       session.CreatedAt,
		UpdatedAt:       session.UpdatedAt,
	}
}

func sessionDTOsFromDomains(sessions []domain.PomodoroSession) []SessionDTOResponse {
	result := make([]SessionDTOResponse, len(sessions))
	for i, session := range sessions {
		result[i] = sessionDTOFromDomain(session)
	}
	return result
}
