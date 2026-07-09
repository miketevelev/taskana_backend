package tasks_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type TaskService struct {
	taskRepository TaskRepository
}

type TaskRepository interface {
	GetTask(
		ctx context.Context,
		userID uuid.UUID,
		taskID uuid.UUID,
	) (domain.Task, error)

	GetTasks(
		ctx context.Context,
		userID uuid.UUID,
		limit *int,
		offset *int,
	) ([]domain.Task, error)

	CreateTask(
		ctx context.Context,
		userID uuid.UUID,
		task domain.Task,
	) (domain.Task, error)

	ChangePosition(
		ctx context.Context,
		userID uuid.UUID,
		task domain.Task,
		oldPosition int,
	) (domain.Task, error)

	PatchTask(
		ctx context.Context,
		userID uuid.UUID,
		task domain.Task,
	) (domain.Task, error)

	PatchTaskWithProjectChange(
		ctx context.Context,
		userID uuid.UUID,
		task domain.Task,
		oldPosition int,
		oldProjectID *uuid.UUID,
	) (domain.Task, error)

	DeleteTask(
		ctx context.Context,
		userID uuid.UUID,
		taskID uuid.UUID,
	) error
}

func NewTaskService(
	taskRepository TaskRepository,
) *TaskService {
	return &TaskService{
		taskRepository: taskRepository,
	}
}
