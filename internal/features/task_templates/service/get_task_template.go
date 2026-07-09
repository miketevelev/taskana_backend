package task_templates_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *TaskTemplatesService) GetTaskTemplate(
	ctx context.Context,
	userID uuid.UUID,
	templateID uuid.UUID,
) (domain.TaskTemplate, error) {
	taskTemplate, err := s.taskTemplateRepository.GetTaskTemplate(
		ctx, userID, templateID,
	)

	if err != nil {
		return domain.TaskTemplate{}, fmt.Errorf(
			"error getting task template: %w", err,
		)
	}

	return taskTemplate, nil
}
