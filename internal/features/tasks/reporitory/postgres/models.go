package tasks_postgres_repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type TaskModel struct {
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

func taskDomainFromModel(taskModel TaskModel) domain.Task {
	return domain.Task{
		ID:            taskModel.ID,
		Version:       taskModel.Version,
		UserID:        taskModel.UserID,
		ProjectID:     taskModel.ProjectID,
		HeadingID:     taskModel.HeadingID,
		TemplateID:    taskModel.TemplateID,
		Title:         taskModel.Title,
		Notes:         taskModel.Notes,
		Status:        taskModel.Status,
		Bucket:        taskModel.Bucket,
		StartDate:     taskModel.StartDate,
		Deadline:      taskModel.Deadline,
		Position:      taskModel.Position,
		IsTimeTracked: taskModel.IsTimeTracked,
		CompletedAt:   taskModel.CompletedAt,
		CreatedAt:     taskModel.CreatedAt,
		UpdatedAt:     taskModel.UpdatedAt,
	}
}

func taskDomainsFromModels(tasks []TaskModel) []domain.Task {
	if len(tasks) == 0 {
		return []domain.Task{}
	}
	taskDomains := make([]domain.Task, len(tasks))

	for i, task := range tasks {
		taskDomains[i] = domain.Task{
			ID:            task.ID,
			Version:       task.Version,
			UserID:        task.UserID,
			ProjectID:     task.ProjectID,
			HeadingID:     task.HeadingID,
			TemplateID:    task.TemplateID,
			Title:         task.Title,
			Notes:         task.Notes,
			Status:        task.Status,
			Bucket:        task.Bucket,
			StartDate:     task.StartDate,
			Deadline:      task.Deadline,
			Position:      task.Position,
			IsTimeTracked: task.IsTimeTracked,
			CompletedAt:   task.CompletedAt,
			CreatedAt:     task.CreatedAt,
			UpdatedAt:     task.UpdatedAt,
		}
	}

	return taskDomains
}

func scanTask(row interface{ Scan(dest ...any) error }) (TaskModel, error) {
	var taskModel TaskModel
	err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.UserID,
		&taskModel.ProjectID,
		&taskModel.HeadingID,
		&taskModel.TemplateID,
		&taskModel.Title,
		&taskModel.Notes,
		&taskModel.Status,
		&taskModel.Bucket,
		&taskModel.StartDate,
		&taskModel.Deadline,
		&taskModel.Position,
		&taskModel.IsTimeTracked,
		&taskModel.EstimatedPomodoros,
		&taskModel.CompletedAt,
		&taskModel.CreatedAt,
		&taskModel.UpdatedAt,
	)
	if err != nil {
		return TaskModel{}, fmt.Errorf("scan row: %w", err)
	}
	return taskModel, nil
}
