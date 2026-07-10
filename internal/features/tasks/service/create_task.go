package tasks_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *TasksService) CreateTask(
	ctx context.Context,
	userID uuid.UUID,
	task domain.Task,
) (domain.Task, error) {
	if err := task.Validate(); err != nil {
		return domain.Task{}, fmt.Errorf(
			"task validation error: %w", err,
		)
	}

	createdTask, err := s.tasksRepository.CreateTask(
		ctx, userID, task,
	)
	if err != nil {
		return domain.Task{}, fmt.Errorf(
			"failed to create task: %w",
			err,
		)
	}

	return createdTask, nil
}
