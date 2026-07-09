package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
)

type RecurrenceType string

const (
	RecurrenceTypeFixed          RecurrenceType = "fixed"
	RecurrenceTypeFromCompletion RecurrenceType = "from_completion"
)

type TargetBucket string

const (
	TargetBucketToday TargetBucket = "today"
	TargetBucketInbox TargetBucket = "inbox"
)

type TaskTemplate struct {
	ID      uuid.UUID `json:"id"`
	Version int       `json:"version"`

	UserID             uuid.UUID      `json:"user_id"`
	ProjectID          *uuid.UUID     `json:"project_id,omitempty"`
	HeadingID          *uuid.UUID     `json:"heading_id,omitempty"`
	Title              string         `json:"title"`
	Notes              *string        `json:"notes,omitempty"`
	RecurrenceRule     string         `json:"recurrence_rule"`
	RecurrenceType     RecurrenceType `json:"recurrence_type"`
	TargetBucket       TargetBucket   `json:"target_bucket"`
	NextExecutionDate  time.Time      `json:"next_execution_date"`
	IsTimeTracked      bool           `json:"is_time_tracked"`
	EstimatedPomodoros int            `json:"estimated_pomodoros"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

func NewTaskTemplate(
	id uuid.UUID,
	version int,
	userID uuid.UUID,
	projectID *uuid.UUID,
	headingID *uuid.UUID,
	title string,
	notes *string,
	recurrenceRule string,
	recurrenceType RecurrenceType,
	targetBucket TargetBucket,
	nextExecutionDate time.Time,
	isTimeTracked bool,
	estimatedPomodoros int,
	createdAt time.Time,
	updatedAt time.Time,
) TaskTemplate {
	return TaskTemplate{
		ID:                 id,
		Version:            version,
		UserID:             userID,
		ProjectID:          projectID,
		HeadingID:          headingID,
		Title:              title,
		Notes:              notes,
		RecurrenceRule:     recurrenceRule,
		RecurrenceType:     recurrenceType,
		TargetBucket:       targetBucket,
		NextExecutionDate:  nextExecutionDate,
		IsTimeTracked:      isTimeTracked,
		EstimatedPomodoros: estimatedPomodoros,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
	}
}

func NewTaskTemplateUninitialized(
	userId uuid.UUID,
	projectID *uuid.UUID,
	headingID *uuid.UUID,
	title string,
	notes *string,
	recurrenceRule string,
	recurrenceType RecurrenceType,
	targetBucket TargetBucket,
	nextExecutionDate time.Time,
	isTimeTracked bool,
	estimatedPomodoros int,
) TaskTemplate {
	now := time.Now().UTC()
	return NewTaskTemplate(
		UninitializedID,
		UninitializedVersion,
		userId,
		projectID,
		headingID,
		title,
		notes,
		recurrenceRule,
		recurrenceType,
		targetBucket,
		nextExecutionDate,
		isTimeTracked,
		estimatedPomodoros,
		now,
		now,
	)
}

func (t *TaskTemplate) Validate() error {
	titleLength := len([]rune(strings.TrimSpace(t.Title)))
	if titleLength < 3 || titleLength > 255 {
		return fmt.Errorf(
			"title must be between 3 and 255 characters: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if t.UserID == uuid.Nil {
		return fmt.Errorf(
			"user_id cannot be empty: %w", core_errors.ErrInvalidArgument,
		)
	}

	if t.ProjectID != nil && *t.ProjectID == uuid.Nil {
		return fmt.Errorf(
			"project_id cannot be an empty UUID: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if t.HeadingID != nil && *t.HeadingID == uuid.Nil {
		return fmt.Errorf(
			"heading_id cannot be an empty UUID: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if t.HeadingID != nil && t.ProjectID == nil {
		return fmt.Errorf(
			"template cannot have a heading without a project: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if strings.TrimSpace(t.RecurrenceRule) == "" {
		return fmt.Errorf(
			"recurrence_rule cannot be empty: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	switch t.RecurrenceType {
	case RecurrenceTypeFixed, RecurrenceTypeFromCompletion:
	default:
		return fmt.Errorf(
			"invalid recurrence_type '%s': %w", t.RecurrenceType,
			core_errors.ErrInvalidArgument,
		)
	}

	switch t.TargetBucket {
	case TargetBucketToday, TargetBucketInbox:
	default:
		return fmt.Errorf(
			"invalid target_bucket '%s': %w", t.TargetBucket,
			core_errors.ErrInvalidArgument,
		)
	}

	if t.NextExecutionDate.IsZero() {
		return fmt.Errorf(
			"next_execution_date cannot be zero: %w",
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

	if t.EstimatedPomodoros < 0 {
		return fmt.Errorf(
			"estimated_pomodoros cannot be negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if t.Notes != nil && len([]rune(strings.TrimSpace(*t.Notes))) > 2000 {
		return fmt.Errorf(
			"notes cannot exceed 2000 characters: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if t.Version < 0 {
		t.Version = 0
	}

	return nil
}

func (t *TaskTemplate) ApplyPatch(patch TaskTemplatePatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate template patch: %w", err)
	}

	tmp := *t

	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}
	if patch.RecurrenceRule.Set {
		tmp.RecurrenceRule = *patch.RecurrenceRule.Value
	}
	if patch.RecurrenceType.Set {
		tmp.RecurrenceType = *patch.RecurrenceType.Value
	}
	if patch.TargetBucket.Set {
		tmp.TargetBucket = *patch.TargetBucket.Value
	}
	if patch.NextExecutionDate.Set {
		tmp.NextExecutionDate = *patch.NextExecutionDate.Value
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

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate template after patch: %w", err)
	}

	*t = tmp
	return nil
}

type TaskTemplatePatch struct {
	ProjectID          Nullable[*uuid.UUID]
	HeadingID          Nullable[*uuid.UUID]
	Title              Nullable[string]
	Notes              Nullable[*string]
	RecurrenceRule     Nullable[string]
	RecurrenceType     Nullable[RecurrenceType]
	TargetBucket       Nullable[TargetBucket]
	NextExecutionDate  Nullable[time.Time]
	IsTimeTracked      Nullable[bool]
	EstimatedPomodoros Nullable[int]
}

func (p *TaskTemplatePatch) Validate() error {
	if p.Title.Set {
		if p.Title.Value == nil || len([]rune(strings.TrimSpace(*p.Title.Value))) < 3 {
			return fmt.Errorf(
				"invalid title in patch (cannot be null or too short): %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if p.RecurrenceRule.Set {
		if p.RecurrenceRule.Value == nil || strings.TrimSpace(*p.RecurrenceRule.Value) == "" {
			return fmt.Errorf(
				"recurrence_rule cannot be null or empty in patch: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if p.NextExecutionDate.Set && p.NextExecutionDate.Value == nil {
		return fmt.Errorf(
			"next_execution_date cannot be null in patch: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.IsTimeTracked.Set && p.IsTimeTracked.Value == nil {
		return fmt.Errorf(
			"is_time_tracked cannot be null in patch: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.RecurrenceType.Set {
		if p.RecurrenceType.Value == nil {
			return fmt.Errorf(
				"recurrence_type cannot be null in patch: %w",
				core_errors.ErrInvalidArgument,
			)
		}
		switch *p.RecurrenceType.Value {
		case RecurrenceTypeFixed, RecurrenceTypeFromCompletion:
		default:
			return fmt.Errorf(
				"invalid patch recurrence type: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if p.TargetBucket.Set {
		if p.TargetBucket.Value == nil {
			return fmt.Errorf(
				"target_bucket cannot be null in patch: %w",
				core_errors.ErrInvalidArgument,
			)
		}
		switch *p.TargetBucket.Value {
		case TargetBucketToday, TargetBucketInbox:
		default:
			return fmt.Errorf(
				"invalid patch target bucket: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if p.EstimatedPomodoros.Set {
		if p.EstimatedPomodoros.Value == nil || *p.EstimatedPomodoros.Value < 0 {
			return fmt.Errorf(
				"patch pomodoros cannot be null or negative: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}
