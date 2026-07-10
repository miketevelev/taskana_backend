package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
)

type Checklist struct {
	ID          uuid.UUID `json:"id"`
	Version     int       `json:"version"`
	UserID      uuid.UUID `json:"user_id"`
	TaskID      uuid.UUID `json:"task_id"`
	Title       string    `json:"title"`
	IsCompleted bool      `json:"is_completed"`
	Position    int       `json:"position"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewChecklist(
	id uuid.UUID,
	version int,
	userID uuid.UUID,
	taskID uuid.UUID,
	title string,
	isCompleted bool,
	position int,
	createdAt time.Time,
	updatedAt time.Time,
) Checklist {
	return Checklist{
		ID:          id,
		Version:     version,
		UserID:      userID,
		TaskID:      taskID,
		Title:       title,
		IsCompleted: isCompleted,
		Position:    position,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func NewChecklistUninitialized(
	userID uuid.UUID,
	taskID uuid.UUID,
	title string,
) Checklist {
	now := time.Now().UTC()
	return NewChecklist(
		UninitializedID,
		UninitializedVersion,
		userID,
		taskID,
		title,
		false,
		1,
		now,
		now,
	)
}

func (c *Checklist) Validate() error {
	titleLength := len([]rune(strings.TrimSpace(c.Title)))
	if titleLength < 3 || titleLength > 255 {
		return fmt.Errorf(
			"title must be between 3 and 255 characters long: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if c.Position < 1 {
		return fmt.Errorf(
			"position must be 1 or greater: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if c.UserID == uuid.Nil {
		return fmt.Errorf(
			"user_id cannot be empty: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if c.TaskID == uuid.Nil {
		return fmt.Errorf(
			"task_id cannot be empty: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if c.CreatedAt.IsZero() || c.UpdatedAt.IsZero() {
		return fmt.Errorf(
			"timestamps cannot be zero: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	if c.UpdatedAt.Before(c.CreatedAt) {
		return fmt.Errorf(
			"updated_at cannot be before created_at: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if c.Version < 0 {
		c.Version = 0
	}

	return nil
}

func (c *Checklist) ApplyPatch(patch ChecklistPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate checklists patch: %w", err)
	}

	tmp := *c

	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}

	if patch.IsCompleted.Set {
		tmp.IsCompleted = *patch.IsCompleted.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate checklists after patch: %w", err)
	}

	*c = tmp

	return nil
}

type ChecklistPatch struct {
	Title       Nullable[string]
	IsCompleted Nullable[bool]
}

func NewChecklistPatch(
	title Nullable[string],
	isCompleted Nullable[bool],
) ChecklistPatch {
	return ChecklistPatch{
		Title:       title,
		IsCompleted: isCompleted,
	}
}

func (p *ChecklistPatch) Validate() error {
	if p.Title.Set && p.Title.Value == nil {
		return fmt.Errorf(
			"'Title' can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.IsCompleted.Set && p.IsCompleted.Value == nil {
		return fmt.Errorf(
			"'IsCompleted' can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}
