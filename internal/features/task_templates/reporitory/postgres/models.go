package task_templates_postgres_repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type TaskTemplateModel struct {
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

func taskTemplateDomainFromModel(taskTemplateModel TaskTemplateModel) domain.TaskTemplate {
	return domain.TaskTemplate{
		ID:                taskTemplateModel.ID,
		Version:           taskTemplateModel.Version,
		UserID:            taskTemplateModel.UserID,
		ProjectID:         taskTemplateModel.ProjectID,
		HeadingID:         taskTemplateModel.HeadingID,
		Title:             taskTemplateModel.Title,
		Notes:             taskTemplateModel.Notes,
		RecurrenceRule:    taskTemplateModel.RecurrenceRule,
		RecurrenceType:    taskTemplateModel.RecurrenceType,
		TargetBucket:      taskTemplateModel.TargetBucket,
		NextExecutionDate: taskTemplateModel.NextExecutionDate,
		IsTimeTracked:     taskTemplateModel.IsTimeTracked,
		CreatedAt:         taskTemplateModel.CreatedAt,
		UpdatedAt:         taskTemplateModel.UpdatedAt,
	}
}

func taskTemplateDomainsFromModels(taskTemplates []TaskTemplateModel) []domain.TaskTemplate {
	if len(taskTemplates) == 0 {
		return []domain.TaskTemplate{}
	}
	taskTemplateDomains := make([]domain.TaskTemplate, len(taskTemplates))

	for i, taskTemplate := range taskTemplates {
		taskTemplateDomains[i] = domain.TaskTemplate{
			ID:                taskTemplate.ID,
			Version:           taskTemplate.Version,
			UserID:            taskTemplate.UserID,
			ProjectID:         taskTemplate.ProjectID,
			HeadingID:         taskTemplate.HeadingID,
			Title:             taskTemplate.Title,
			Notes:             taskTemplate.Notes,
			RecurrenceRule:    taskTemplate.RecurrenceRule,
			RecurrenceType:    taskTemplate.RecurrenceType,
			TargetBucket:      taskTemplate.TargetBucket,
			NextExecutionDate: taskTemplate.NextExecutionDate,
			IsTimeTracked:     taskTemplate.IsTimeTracked,
			CreatedAt:         taskTemplate.CreatedAt,
			UpdatedAt:         taskTemplate.UpdatedAt,
		}
	}

	return taskTemplateDomains
}

func scanTaskTemplate(row interface{ Scan(dest ...any) error }) (
	TaskTemplateModel,
	error,
) {
	var taskTemplateModel TaskTemplateModel
	err := row.Scan(
		&taskTemplateModel.ID,
		&taskTemplateModel.Version,
		&taskTemplateModel.UserID,
		&taskTemplateModel.ProjectID,
		&taskTemplateModel.HeadingID,
		&taskTemplateModel.Title,
		&taskTemplateModel.Notes,
		&taskTemplateModel.RecurrenceRule,
		&taskTemplateModel.RecurrenceType,
		&taskTemplateModel.TargetBucket,
		&taskTemplateModel.NextExecutionDate,
		&taskTemplateModel.IsTimeTracked,
		&taskTemplateModel.CreatedAt,
		&taskTemplateModel.UpdatedAt,
	)
	if err != nil {
		return TaskTemplateModel{}, fmt.Errorf("scan task template: %w", err)
	}
	return taskTemplateModel, nil
}
