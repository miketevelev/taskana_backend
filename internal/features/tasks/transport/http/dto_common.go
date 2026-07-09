package tasks_transport_http

import (
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type TaskDTOResponse struct {
	ID                 uuid.UUID         `json:"id"`
	Version            int               `json:"version"`
	UserID             uuid.UUID         `json:"user_id"`
	ProjectID          *uuid.UUID        `json:"project_id,omitempty"`
	HeadingID          *uuid.UUID        `json:"heading_id,omitempty"`
	TemplateID         *uuid.UUID        `json:"template_id,omitempty"`
	Title              string            `json:"title"`
	Notes              *string           `json:"notes,omitempty"`
	Status             domain.TaskStatus `json:"status"`
	Bucket             domain.TaskBucket `json:"bucket"`
	StartDate          *time.Time        `json:"start_date,omitempty"`
	Deadline           *time.Time        `json:"deadline,omitempty"`
	Position           int               `json:"position"`
	IsTimeTracked      bool              `json:"is_time_tracked"`
	EstimatedPomodoros int               `json:"estimated_pomodoros"`
	CompletedAt        *time.Time        `json:"completed_at,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}

func taskDTOFromDomain(task domain.Task) TaskDTOResponse {
	return TaskDTOResponse{
		ID:                 task.ID,
		Version:            task.Version,
		UserID:             task.UserID,
		ProjectID:          task.ProjectID,
		HeadingID:          task.HeadingID,
		TemplateID:         task.TemplateID,
		Title:              task.Title,
		Notes:              task.Notes,
		Status:             task.Status,
		Bucket:             task.Bucket,
		StartDate:          task.StartDate,
		Deadline:           task.Deadline,
		Position:           task.Position,
		IsTimeTracked:      task.IsTimeTracked,
		EstimatedPomodoros: task.EstimatedPomodoros,
		CompletedAt:        task.CompletedAt,
		CreatedAt:          task.CreatedAt,
		UpdatedAt:          task.UpdatedAt,
	}
}

func taskDTOsFromDomains(tasks []domain.Task) []TaskDTOResponse {
	result := make([]TaskDTOResponse, len(tasks))
	for i, task := range tasks {
		result[i] = taskDTOFromDomain(task)
	}
	return result
}
