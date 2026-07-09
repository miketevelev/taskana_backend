package task_templates_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *TaskTemplatesService) DeleteTaskTemplate(
	ctx context.Context,
	userID uuid.UUID,
	taskTemplateID uuid.UUID,
) error {
	if err := s.taskTemplateRepository.DeleteTaskTemplate(
		ctx, userID, taskTemplateID,
	); err != nil {
		return fmt.Errorf(
			"delete task template: %w", err,
		)
	}

	return nil
}
