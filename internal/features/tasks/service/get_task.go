package tasks_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *TaskService) GetTask(
	ctx context.Context,
	userID uuid.UUID,
	taskID uuid.UUID,
) (domain.Task, error) {
	task, err := s.taskRepository.GetTask(ctx, userID, taskID)
	if err != nil {
		return domain.Task{}, fmt.Errorf(
			"error getting task %s: %w", taskID, err,
		)
	}

	return task, nil
}
