package task_templates_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
)

type TaskTemplatesService struct {
	taskTemplateRepository TaskTemplateRepository
}

type TaskTemplateRepository interface {
	CreateTaskTemplate(
		ctx context.Context,
		userID uuid.UUID,
		taskTemplate domain.TaskTemplate,
	) (domain.TaskTemplate, error)
}

func NewTaskTemplatesService(
	taskTemplateRepository TaskTemplateRepository,
) *TaskTemplatesService {
	return &TaskTemplatesService{
		taskTemplateRepository: taskTemplateRepository,
	}
}

func ParseRecurrenceType(value string) (domain.RecurrenceType, error) {
	switch domain.RecurrenceType(value) {
	case domain.RecurrenceTypeFixed, domain.RecurrenceTypeFromCompletion:
		return domain.RecurrenceType(value), nil
	default:
		return "", fmt.Errorf(
			"invalid recurrence_type %q: %w", value,
			core_errors.ErrInvalidArgument,
		)
	}
}

func ParseTargetBucket(value string) (domain.TargetBucket, error) {
	switch domain.TargetBucket(value) {
	case domain.TargetBucketToday, domain.TargetBucketInbox:
		return domain.TargetBucket(value), nil
	default:
		return "", fmt.Errorf(
			"invalid target_bucket %q: %w", value,
			core_errors.ErrInvalidArgument,
		)
	}
}
