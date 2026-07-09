package tasks_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *TaskService) DeleteTask(
	ctx context.Context,
	userID uuid.UUID,
	taskID uuid.UUID,
) error {
	if err := s.taskRepository.DeleteTask(
		ctx, userID, taskID,
	); err != nil {
		return fmt.Errorf(
			"delete task: %w", err,
		)
	}

	return nil
}
