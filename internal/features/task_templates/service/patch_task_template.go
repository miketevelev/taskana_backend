package task_templates_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *TaskTemplatesService) PatchTaskTemplate(
	ctx context.Context,
	userID uuid.UUID,
	taskTemplateID uuid.UUID,
	patch domain.TaskTemplatePatch,
) (domain.TaskTemplate, error) {
	taskTemplate, err := s.taskTemplateRepository.GetTaskTemplate(ctx, userID, taskTemplateID)
	if err != nil {
		return domain.TaskTemplate{}, fmt.Errorf(
			"error while fetching task template: %w", err,
		)
	}

	if err := taskTemplate.ApplyPatch(patch); err != nil {
		return domain.TaskTemplate{}, fmt.Errorf(
			"error while applying patch to task template: %w",
			err,
		)
	}

	patchedTaskTemplate, err := s.taskTemplateRepository.PatchTaskTemplate(
		ctx, userID, taskTemplate,
	)
	if err != nil {
		return domain.TaskTemplate{}, fmt.Errorf(
			"error while saving patched task template: %w", err,
		)
	}

	return patchedTaskTemplate, nil
}
