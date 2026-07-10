package tasks_service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_recurrence "github.com/miketevelev/taskana_backend/internal/core/recurrence"
)

type TasksService struct {
	tasksRepository TasksRepository
}

type TasksRepository interface {
	GetTask(
		ctx context.Context,
		userID uuid.UUID,
		taskID uuid.UUID,
	) (domain.Task, error)

	GetTasks(
		ctx context.Context,
		userID uuid.UUID,
		limit *int,
		offset *int,
	) ([]domain.Task, error)

	CreateTask(
		ctx context.Context,
		userID uuid.UUID,
		task domain.Task,
	) (domain.Task, error)

	ChangePosition(
		ctx context.Context,
		userID uuid.UUID,
		task domain.Task,
		oldPosition int,
	) (domain.Task, error)

	PatchTask(
		ctx context.Context,
		userID uuid.UUID,
		task domain.Task,
	) (domain.Task, error)

	PatchTaskWithProjectChange(
		ctx context.Context,
		userID uuid.UUID,
		task domain.Task,
		oldPosition int,
		oldProjectID *uuid.UUID,
	) (domain.Task, error)

	DeleteTask(
		ctx context.Context,
		userID uuid.UUID,
		taskID uuid.UUID,
	) error

	UpdateTemplateNextExecution(
		ctx context.Context,
		userID, templateID uuid.UUID,
		nextDate time.Time,
	) error

	ListDueFixedTemplates(
		ctx context.Context,
		asOf time.Time,
	) ([]domain.TaskTemplate, error)

	ProcessNextDueFixedTemplateTx(
		ctx context.Context,
		asOf time.Time,
		processFn func(template domain.TaskTemplate) (
			domain.Task,
			time.Time,
			error,
		),
	) (bool, error)
}

func NewTaskService(
	taskRepository TasksRepository,
) *TasksService {
	return &TasksService{
		tasksRepository: taskRepository,
	}
}

func (s *TasksService) ProcessFixedRecurrences(
	ctx context.Context,
	asOf time.Time,
) error {
	for {
		processed, err := s.tasksRepository.ProcessNextDueFixedTemplateTx(
			ctx,
			asOf,
			s.generateFixedTaskData,
		)

		if err != nil {
			return err
		}

		if !processed {
			break
		}
	}

	return nil
}

func (s *TasksService) generateFixedTaskData(
	template domain.TaskTemplate,
) (domain.Task, time.Time, error) {
	bucket := domain.TaskBucketInbox
	if template.TargetBucket == domain.TargetBucketToday {
		bucket = domain.TaskBucketToday
	}

	task := domain.NewTaskUninitialized(
		template.UserID,
		template.ProjectID,
		template.HeadingID,
		&template.ID,
		template.Title,
		template.Notes,
		bucket,
		nil,
		nil,
		template.IsTimeTracked,
		template.EstimatedPomodoros,
	)

	task.ID = uuid.New()
	task.Version = 1

	if err := task.Validate(); err != nil {
		return domain.Task{}, time.Time{}, fmt.Errorf(
			"invalid generated task: %w", err,
		)
	}

	nextDate, err := core_recurrence.NextFixedDate(
		template.RecurrenceRule,
		template.NextExecutionDate,
	)
	if err != nil {
		return domain.Task{}, time.Time{}, fmt.Errorf(
			"calculate next fixed date: %w", err,
		)
	}

	return task, nextDate, nil
}

//func (s *TasksService) generateFixedTask(
//	ctx context.Context,
//	template domain.TaskTemplate,
//) error {
//	bucket := domain.TaskBucketInbox
//	if template.TargetBucket == domain.TargetBucketToday {
//		bucket = domain.TaskBucketToday
//	}
//
//	task := domain.NewTaskUninitialized(
//		template.UserID,
//		template.ProjectID,
//		template.HeadingID,
//		&template.ID,
//		template.Title,
//		template.Notes,
//		bucket,
//		nil,
//		nil,
//		template.IsTimeTracked,
//		template.EstimatedPomodoros,
//	)
//
//	if err := task.Validate(); err != nil {
//		return fmt.Errorf("invalid generated task: %w", err)
//	}
//
//	nextDate, err := core_recurrence.NextFixedDate(
//		template.RecurrenceRule,
//		template.NextExecutionDate,
//	)
//	if err != nil {
//		return fmt.Errorf("calculate next fixed date: %w", err)
//	}
//
//	if _, err := s.tasksRepository.CreateTask(
//		ctx, template.UserID, task,
//	); err != nil {
//		return fmt.Errorf("create fixed recurring task: %w", err)
//	}
//
//	if err := s.tasksRepository.UpdateTemplateNextExecution(
//		ctx, template.UserID, template.ID, nextDate,
//	); err != nil {
//		return fmt.Errorf("update template next execution date: %w", err)
//	}
//
//	return nil
//}
