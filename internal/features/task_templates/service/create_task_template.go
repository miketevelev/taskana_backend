package task_templates_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

func (s *TaskTemplatesService) CreateTaskTemplate(
	ctx context.Context,
	userID uuid.UUID,
	taskTemplate domain.TaskTemplate,
) (domain.TaskTemplate, error) {
	if taskTemplate.ID == uuid.Nil {
		taskTemplate.ID = uuid.New()
	}

	if taskTemplate.Version == -1 {
		taskTemplate.Version = 1
	}

	if err := taskTemplate.Validate(); err != nil {
		return domain.TaskTemplate{},
			fmt.Errorf("task template validation failed: %w", err)
	}

	createdTaskTemplate, err := s.taskTemplateRepository.CreateTaskTemplate(
		ctx, userID, taskTemplate,
	)
	if err != nil {
		return domain.TaskTemplate{}, fmt.Errorf(
			"create task template failed: %w", err,
		)
	}

	return createdTaskTemplate, nil
}
