package tasks_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *TasksService) DeleteTask(
	ctx context.Context,
	userID uuid.UUID,
	taskID uuid.UUID,
) error {
	if err := s.tasksRepository.DeleteTask(
		ctx, userID, taskID,
	); err != nil {
		return fmt.Errorf(
			"delete task: %w", err,
		)
	}

	return nil
}
