package tasks_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *TasksService) ChangePosition(
	ctx context.Context,
	userID uuid.UUID,
	taskID uuid.UUID,
	newPosition int,
) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, userID, taskID)
	if err != nil {
		return domain.Task{}, fmt.Errorf(
			"error getting task: %w", err,
		)
	}

	oldPosition := task.Position

	if oldPosition == newPosition {
		return task, nil
	}

	task.Position = newPosition

	updatedTask, err := s.tasksRepository.ChangePosition(
		ctx, userID, task, oldPosition,
	)
	if err != nil {
		return domain.Task{}, fmt.Errorf(
			"error changing position in repository: %w", err,
		)
	}

	return updatedTask, nil
}
