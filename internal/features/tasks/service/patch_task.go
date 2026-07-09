package tasks_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *TaskService) PatchTask(
	ctx context.Context,
	userID uuid.UUID,
	taskID uuid.UUID,
	patch domain.TaskPatch,
) (domain.Task, error) {
	task, err := s.taskRepository.GetTask(ctx, userID, taskID)
	if err != nil {
		return domain.Task{}, fmt.Errorf(
			"error while fetching task: %w", err,
		)
	}

	oldProjectID := task.ProjectID
	oldPosition := task.Position

	if err := task.ApplyPatch(patch); err != nil {
		return domain.Task{}, fmt.Errorf(
			"error while applying patch to task: %w", err,
		)
	}

	projectChanged := false
	if patch.ProjectID.Set {
		if (oldProjectID == nil && task.ProjectID != nil) ||
			(oldProjectID != nil && task.ProjectID == nil) ||
			(oldProjectID != nil && task.ProjectID != nil && *oldProjectID != *task.ProjectID) {
			projectChanged = true
		}
	}

	var patchedTask domain.Task

	if projectChanged {
		patchedTask, err = s.taskRepository.PatchTaskWithProjectChange(
			ctx, userID, task, oldPosition, oldProjectID,
		)
	} else {
		patchedTask, err = s.taskRepository.PatchTask(
			ctx, userID, task,
		)
	}

	if err != nil {
		return domain.Task{}, fmt.Errorf(
			"error while saving patched task: %w", err,
		)
	}

	return patchedTask, nil
}
