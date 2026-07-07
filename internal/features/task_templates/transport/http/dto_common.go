package task_templates_transport_http

import (
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type TaskTemplateDTOResponse struct {
	ID                 uuid.UUID             `json:"id"`
	Version            int                   `json:"version"`
	UserID             uuid.UUID             `json:"user_id"`
	ProjectID          *uuid.UUID            `json:"project_id,omitempty"`
	HeadingID          *uuid.UUID            `json:"heading_id,omitempty"`
	Title              string                `json:"title"`
	Notes              *string               `json:"notes,omitempty"`
	RecurrenceRule     string                `json:"recurrence_rule"`
	RecurrenceType     domain.RecurrenceType `json:"recurrence_type"`
	TargetBucket       domain.TargetBucket   `json:"target_bucket"`
	NextExecutionDate  time.Time             `json:"next_execution_date"`
	IsTimeTracked      bool                  `json:"is_time_tracked"`
	EstimatedPomodoros int                   `json:"estimated_pomodoros"`
	CreatedAt          time.Time             `json:"created_at"`
	UpdatedAt          time.Time             `json:"updated_at"`
}

func taskTemplateDTOFromDomain(taskTemplate domain.TaskTemplate) TaskTemplateDTOResponse {
	return TaskTemplateDTOResponse{
		ID:                 taskTemplate.ID,
		Version:            taskTemplate.Version,
		UserID:             taskTemplate.UserID,
		ProjectID:          taskTemplate.ProjectID,
		HeadingID:          taskTemplate.HeadingID,
		Title:              taskTemplate.Title,
		Notes:              taskTemplate.Notes,
		RecurrenceRule:     taskTemplate.RecurrenceRule,
		RecurrenceType:     taskTemplate.RecurrenceType,
		TargetBucket:       taskTemplate.TargetBucket,
		NextExecutionDate:  taskTemplate.NextExecutionDate,
		IsTimeTracked:      taskTemplate.IsTimeTracked,
		EstimatedPomodoros: taskTemplate.EstimatedPomodoros,
		CreatedAt:          taskTemplate.CreatedAt,
		UpdatedAt:          taskTemplate.UpdatedAt,
	}
}

func taskTemplateDTOsFromDomains(
	taskTemplates []domain.
		TaskTemplate,
) []TaskTemplateDTOResponse {
	result := make([]TaskTemplateDTOResponse, len(taskTemplates))
	for i, t := range taskTemplates {
		result[i] = taskTemplateDTOFromDomain(t)
	}
	return result
}
