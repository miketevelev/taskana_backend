package task_templates_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
)

func (s *TaskTemplatesService) GetTaskTemplates(
	ctx context.Context,
	userId uuid.UUID,
	limit *int,
	offset *int,
) ([]domain.TaskTemplate, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf(
			"limit must be non-negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf(
			"offset must be non-negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	taskTemplates, err := s.taskTemplateRepository.GetTaskTemplates(
		ctx, userId, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("get task templates from repository: %w", err)
	}

	return taskTemplates, nil
}
