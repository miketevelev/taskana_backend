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
}

func NewTaskService(
	taskRepository TaskRepository,
) *TaskService {
	return &TaskService{
		taskRepository: taskRepository,
	}
}
