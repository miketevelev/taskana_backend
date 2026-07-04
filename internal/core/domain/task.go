package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
)

type TaskStatus string

const (
	TaskStatusOpen      TaskStatus = "open"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusCanceled  TaskStatus = "canceled"
)

type TaskBucket string

const (
	TaskBucketInbox   TaskBucket = "inbox"
	TaskBucketToday   TaskBucket = "today"
	TaskBucketAnytime TaskBucket = "anytime"
	TaskBucketSomeday TaskBucket = "someday"
)

type Task struct {
	ID      uuid.UUID `json:"id"`
	Version int       `json:"version"`

	UserID             uuid.UUID  `json:"user_id"`
	ProjectID          *uuid.UUID `json:"project_id,omitempty"`
	HeadingID          *uuid.UUID `json:"heading_id,omitempty"`
	TemplateID         *uuid.UUID `json:"template_id,omitempty"`
	Title              string     `json:"title"`
	Notes              *string    `json:"notes,omitempty"`
	Status             TaskStatus `json:"status"`
	Bucket             TaskBucket `json:"bucket"`
	StartDate          *time.Time `json:"start_date,omitempty"`
	Deadline           *time.Time `json:"deadline,omitempty"`
	Position           int        `json:"position"`
	IsTimeTracked      bool       `json:"is_time_tracked"`
	EstimatedPomodoros int        `json:"estimated_pomodoros"`
	CompletedAt        *time.Time `json:"completed_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

func NewTask(
	id uuid.UUID,
	version int,
	userID uuid.UUID,
	projectID *uuid.UUID,
	headingID *uuid.UUID,
	templateID *uuid.UUID,
	title string,
	notes *string,
	status TaskStatus,
	bucket TaskBucket,
	startDate *time.Time,
	deadline *time.Time,
	position int,
	isTimeTracked bool,
	estimatedPomodoros int,
	completedAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) Task {
	return Task{
		ID:                 id,
		Version:            version,
		UserID:             userID,
		ProjectID:          projectID,
		HeadingID:          headingID,
		TemplateID:         templateID,
		Title:              title,
		Notes:              notes,
		Status:             status,
		Bucket:             bucket,
		StartDate:          startDate,
		Deadline:           deadline,
		Position:           position,
		IsTimeTracked:      isTimeTracked,
		EstimatedPomodoros: estimatedPomodoros,
		CompletedAt:        completedAt,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
	}
}

func NewTaskUninitialized(
	userID uuid.UUID,
	projectID *uuid.UUID,
	headingID *uuid.UUID,
	title string,
	notes *string,
	bucket TaskBucket,
	startDate *time.Time,
	deadline *time.Time,
	isTimeTracked bool,
	estimatedPomodoros int,
) Task {
	now := time.Now().UTC()
	return Task{
		ID:                 UninitializedID,
		Version:            UninitializedVersion,
		UserID:             userID,
		ProjectID:          projectID,
		HeadingID:          headingID,
		TemplateID:         nil,
		Title:              title,
		Notes:              notes,
		Status:             TaskStatusOpen,
		Bucket:             bucket,
		StartDate:          startDate,
		Deadline:           deadline,
		Position:           1,
		IsTimeTracked:      isTimeTracked,
		EstimatedPomodoros: estimatedPomodoros,
		CompletedAt:        nil,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

func (t *Task) Validate() error {
	titleLength := len([]rune(strings.TrimSpace(t.Title)))
	if titleLength < 3 || titleLength > 255 {
		return fmt.Errorf(
			"title must be between 3 and 255 characters long: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if t.Position < 1 {
		return fmt.Errorf(
			"position must be 1 or greater: %w", core_errors.ErrInvalidArgument,
		)
	}

	if t.UserID == uuid.Nil {
		return fmt.Errorf(
			"user_id cannot be empty: %w", core_errors.ErrInvalidArgument,
		)
	}

	if t.HeadingID != nil && t.ProjectID == nil {
		return fmt.Errorf(
			"task cannot have a heading without a project: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if t.CreatedAt.IsZero() || t.UpdatedAt.IsZero() {
		return fmt.Errorf(
			"timestamps cannot be zero: %w", core_errors.ErrInvalidArgument,
		)
	}

	if t.UpdatedAt.Before(t.CreatedAt) {
		return fmt.Errorf(
			"updated_at cannot be before created_at: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	switch t.Status {
	case TaskStatusOpen, TaskStatusCompleted, TaskStatusCanceled:
	default:
		return fmt.Errorf(
			"invalid task status '%s': %w", t.Status,
			core_errors.ErrInvalidArgument,
		)
	}

	switch t.Bucket {
	case TaskBucketInbox, TaskBucketToday, TaskBucketAnytime, TaskBucketSomeday:
	default:
		return fmt.Errorf(
			"invalid task bucket '%s': %w", t.Bucket,
			core_errors.ErrInvalidArgument,
		)
	}

	isClosed := t.Status == TaskStatusCompleted || t.Status == TaskStatusCanceled
	if isClosed && t.CompletedAt == nil {
		return fmt.Errorf(
			"completed_at cannot be nil if task is %s: %w", t.Status,
			core_errors.ErrInvalidArgument,
		)
	}
	if !isClosed && t.CompletedAt != nil {
		return fmt.Errorf(
			"completed_at must be nil if task is open: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	if t.CompletedAt != nil && t.CompletedAt.Before(t.CreatedAt) {
		return fmt.Errorf(
			"completed_at cannot be before created_at: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if t.StartDate != nil && t.Deadline != nil {
		if t.Deadline.Before(*t.StartDate) {
			return fmt.Errorf(
				"deadline cannot be before start_date: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if t.EstimatedPomodoros < 0 {
		return fmt.Errorf(
			"estimated_pomodoros cannot be negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if t.Notes != nil {
		if len([]rune(strings.TrimSpace(*t.Notes))) > 2000 {
			return fmt.Errorf(
				"notes cannot exceed 2000 characters: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if t.Version < 0 {
		t.Version = 0
	}

	return nil
}

func (t *Task) ApplyPatch(patch TaskPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate task patch: %w", err)
	}

	tmp := *t
	now := time.Now().UTC()

	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}
	if patch.Status.Set {
		tmp.Status = *patch.Status.Value
	}
	if patch.Bucket.Set {
		tmp.Bucket = *patch.Bucket.Value
	}
	if patch.IsTimeTracked.Set {
		tmp.IsTimeTracked = *patch.IsTimeTracked.Value
	}
	if patch.EstimatedPomodoros.Set {
		tmp.EstimatedPomodoros = *patch.EstimatedPomodoros.Value
	}

	if patch.ProjectID.Set {
		if patch.ProjectID.Value == nil {
			tmp.ProjectID = nil
			tmp.HeadingID = nil
		} else {
			tmp.ProjectID = *patch.ProjectID.Value
		}
	}
	if patch.HeadingID.Set {
		if patch.HeadingID.Value == nil {
			tmp.HeadingID = nil
		} else {
			tmp.HeadingID = *patch.HeadingID.Value
		}
	}
	if patch.Notes.Set {
		if patch.Notes.Value == nil {
			tmp.Notes = nil
		} else {
			tmp.Notes = *patch.Notes.Value
		}
	}
	if patch.StartDate.Set {
		if patch.StartDate.Value == nil {
			tmp.StartDate = nil
		} else {
			tmp.StartDate = *patch.StartDate.Value
		}
	}
	if patch.Deadline.Set {
		if patch.Deadline.Value == nil {
			tmp.Deadline = nil
		} else {
			tmp.Deadline = *patch.Deadline.Value
		}
	}
	if patch.CompletedAt.Set {
		if patch.CompletedAt.Value == nil {
			tmp.CompletedAt = nil
		} else {
			tmp.CompletedAt = *patch.CompletedAt.Value
		}
	}

	isClosed := tmp.Status == TaskStatusCompleted || tmp.Status == TaskStatusCanceled
	wasClosed := t.Status == TaskStatusCompleted || t.Status == TaskStatusCanceled

	if patch.Status.Set && isClosed && !wasClosed && !patch.CompletedAt.Set {
		tmp.CompletedAt = &now
	}
	if patch.Status.Set && !isClosed && wasClosed && !patch.CompletedAt.Set {
		tmp.CompletedAt = nil
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate task after patch: %w", err)
	}

	*t = tmp

	return nil
}

type TaskPatch struct {
	ProjectID          Nullable[*uuid.UUID]
	HeadingID          Nullable[*uuid.UUID]
	Title              Nullable[string]
	Notes              Nullable[*string]
	Status             Nullable[TaskStatus]
	Bucket             Nullable[TaskBucket]
	StartDate          Nullable[*time.Time]
	Deadline           Nullable[*time.Time]
	IsTimeTracked      Nullable[bool]
	EstimatedPomodoros Nullable[int]
	CompletedAt        Nullable[*time.Time]
}

func NewTaskPatch(
	projectID Nullable[*uuid.UUID],
	headingID Nullable[*uuid.UUID],
	title Nullable[string],
	notes Nullable[*string],
	status Nullable[TaskStatus],
	bucket Nullable[TaskBucket],
	startDate Nullable[*time.Time],
	deadline Nullable[*time.Time],
	isTimeTracked Nullable[bool],
	estimatedPomodoros Nullable[int],
	completedAt Nullable[*time.Time],
) TaskPatch {
	return TaskPatch{
		ProjectID:          projectID,
		HeadingID:          headingID,
		Title:              title,
		Notes:              notes,
		Status:             status,
		Bucket:             bucket,
		StartDate:          startDate,
		Deadline:           deadline,
		IsTimeTracked:      isTimeTracked,
		EstimatedPomodoros: estimatedPomodoros,
		CompletedAt:        completedAt,
	}
}

func (p *TaskPatch) Validate() error {
	if p.Title.Set && p.Title.Value == nil {
		return fmt.Errorf(
			"'Title' can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	if p.Status.Set && p.Status.Value == nil {
		return fmt.Errorf(
			"'Status' can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	if p.Bucket.Set && p.Bucket.Value == nil {
		return fmt.Errorf(
			"'Bucket' can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	if p.IsTimeTracked.Set && p.IsTimeTracked.Value == nil {
		return fmt.Errorf(
			"'IsTimeTracked' can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	if p.EstimatedPomodoros.Set && p.EstimatedPomodoros.Value == nil {
		return fmt.Errorf(
			"'EstimatedPomodoros' can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.Status.Set && p.Status.Value != nil {
		switch *p.Status.Value {
		case TaskStatusOpen, TaskStatusCompleted, TaskStatusCanceled:
		default:
			return fmt.Errorf(
				"invalid patch task status '%s': %w", *p.Status.Value,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if p.Bucket.Set && p.Bucket.Value != nil {
		switch *p.Bucket.Value {
		case TaskBucketInbox, TaskBucketToday, TaskBucketAnytime, TaskBucketSomeday:
		default:
			return fmt.Errorf(
				"invalid patch task bucket '%s': %w", *p.Bucket.Value,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if p.EstimatedPomodoros.Set && p.EstimatedPomodoros.Value != nil && *p.EstimatedPomodoros.Value < 0 {
		return fmt.Errorf(
			"patch 'EstimatedPomodoros' cannot be negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.ProjectID.Set && p.ProjectID.Value != nil && *p.ProjectID.Value != nil {
		if **p.ProjectID.Value == uuid.Nil {
			return fmt.Errorf(
				"patch 'ProjectID' cannot be an empty UUID: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if p.HeadingID.Set && p.HeadingID.Value != nil && *p.HeadingID.Value != nil {
		if **p.HeadingID.Value == uuid.Nil {
			return fmt.Errorf(
				"patch 'HeadingID' cannot be an empty UUID: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}
