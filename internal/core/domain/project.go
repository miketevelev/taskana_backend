package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
)

type ProjectStatus string

const (
	ProjectStatusActive    ProjectStatus = "active"
	ProjectStatusCompleted ProjectStatus = "completed"
	ProjectStatusDropped   ProjectStatus = "dropped"
)

type Project struct {
	ID          uuid.UUID     `json:"id"`
	Version     int           `json:"version"`
	UserID      uuid.UUID     `json:"user_id"`
	AreaID      *uuid.UUID    `json:"area_id,omitempty"`
	Title       string        `json:"title"`
	Notes       *string       `json:"notes,omitempty"`
	Status      ProjectStatus `json:"status"`
	Position    int           `json:"position"`
	Deadline    *time.Time    `json:"deadline,omitempty"`
	CompletedAt *time.Time    `json:"completed_at,omitempty"`
	CreatedAt   time.Time     `json:"created_at,omitempty"`
	UpdatedAt   time.Time     `json:"updated_at,omitempty"`
}

func NewProject(
	id uuid.UUID,
	version int,
	userID uuid.UUID,
	areaID *uuid.UUID,
	title string,
	notes *string,
	status ProjectStatus,
	position int,
	deadline *time.Time,
	completedAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) Project {
	return Project{
		ID:          id,
		Version:     version,
		UserID:      userID,
		AreaID:      areaID,
		Title:       title,
		Notes:       notes,
		Status:      status,
		Position:    position,
		Deadline:    deadline,
		CompletedAt: completedAt,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func NewProjectUninitialized(
	userID uuid.UUID,
	areaID *uuid.UUID,
	title string,
	notes *string,
	deadline *time.Time,
) Project {
	now := time.Now().UTC()
	return Project{
		ID:          UninitializedID,
		Version:     UninitializedVersion,
		UserID:      userID,
		AreaID:      areaID,
		Title:       title,
		Notes:       notes,
		Status:      ProjectStatusActive,
		Position:    0,
		Deadline:    deadline,
		CompletedAt: nil,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (p *Project) Validate() error {
	titleLength := len([]rune(strings.TrimSpace(p.Title)))
	if titleLength < 3 || titleLength > 100 {
		return fmt.Errorf(
			"title must be between 3 and 100 characters long: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.Position < 1 && p.Position != 0 {
		return fmt.Errorf(
			"position must be 1 or greater: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.CreatedAt.IsZero() || p.UpdatedAt.IsZero() {
		return fmt.Errorf(
			"timestamps cannot be zero: %w", core_errors.ErrInvalidArgument,
		)
	}
	if p.UpdatedAt.Before(p.CreatedAt) {
		return fmt.Errorf(
			"updated_at cannot be before created_at: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	switch p.Status {
	case ProjectStatusActive, ProjectStatusCompleted, ProjectStatusDropped:
	default:
		return fmt.Errorf(
			"invalid project status '%s': %w",
			p.Status, core_errors.ErrInvalidArgument,
		)
	}

	if p.Status == ProjectStatusCompleted && p.CompletedAt == nil {
		return fmt.Errorf(
			"completed_at cannot be nil if project is completed: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	if p.Status != ProjectStatusCompleted && p.CompletedAt != nil {
		return fmt.Errorf(
			"completed_at must be nil if project is not completed: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	if p.CompletedAt != nil && p.CompletedAt.Before(p.CreatedAt) {
		return fmt.Errorf(
			"completed_at cannot be before created_at: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.AreaID != nil && *p.AreaID == uuid.Nil {
		return fmt.Errorf(
			"area_id cannot be empty UUID if provided: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.Notes != nil {
		notesLength := len([]rune(strings.TrimSpace(*p.Notes)))
		if notesLength > 2000 {
			return fmt.Errorf(
				"notes cannot exceed 2000 characters: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if p.Deadline != nil && p.Deadline.IsZero() {
		return fmt.Errorf(
			"deadline cannot be a zero time if provided: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.Version < 0 {
		p.Version = 0
	}

	return nil
}

func (p *Project) ApplyPatch(patch ProjectPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate project patch: %w", err)
	}

	tmp := *p
	now := time.Now().UTC()

	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}
	if patch.Position.Set {
		tmp.Position = *patch.Position.Value
	}
	if patch.Status.Set {
		tmp.Status = *patch.Status.Value
	}
	if patch.AreaID.Set {
		tmp.AreaID = *patch.AreaID.Value
	}
	if patch.Notes.Set {
		tmp.Notes = *patch.Notes.Value
	}
	if patch.Deadline.Set {
		tmp.Deadline = *patch.Deadline.Value
	}
	if patch.CompletedAt.Set {
		tmp.CompletedAt = *patch.CompletedAt.Value
	}

	if patch.Status.Set && tmp.Status == ProjectStatusCompleted && !patch.CompletedAt.Set {
		tmp.CompletedAt = &now
	}
	if patch.Status.Set && tmp.Status != ProjectStatusCompleted && !patch.CompletedAt.Set {
		tmp.CompletedAt = nil
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate project after patch: %w", err)
	}

	*p = tmp

	return nil
}

type ProjectPatch struct {
	AreaID      Nullable[*uuid.UUID]
	Title       Nullable[string]
	Notes       Nullable[*string]
	Status      Nullable[ProjectStatus]
	Position    Nullable[int]
	Deadline    Nullable[*time.Time]
	CompletedAt Nullable[*time.Time]
}

func NewProjectPatch(
	areaID Nullable[*uuid.UUID],
	title Nullable[string],
	notes Nullable[*string],
	status Nullable[ProjectStatus],
	position Nullable[int],
	deadline Nullable[*time.Time],
	completedAt Nullable[*time.Time],
) ProjectPatch {
	return ProjectPatch{
		AreaID:      areaID,
		Title:       title,
		Notes:       notes,
		Status:      status,
		Position:    position,
		Deadline:    deadline,
		CompletedAt: completedAt,
	}
}

func (p *ProjectPatch) Validate() error {
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

	if p.Position.Set {
		if p.Position.Value == nil {
			return fmt.Errorf(
				"'Position' can't be patched to NULL: %w",
				core_errors.ErrInvalidArgument,
			)
		}
		if *p.Position.Value < 1 {
			return fmt.Errorf(
				"'Position' must be 1 or greater: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if p.Status.Set && p.Status.Value != nil {
		status := *p.Status.Value
		switch status {
		case ProjectStatusActive, ProjectStatusCompleted, ProjectStatusDropped:
		default:
			return fmt.Errorf(
				"invalid patch project status '%s': %w",
				status, core_errors.ErrInvalidArgument,
			)
		}
	}

	if p.AreaID.Set && p.AreaID.Value != nil && *p.AreaID.Value != nil {
		if **p.AreaID.Value == uuid.Nil {
			return fmt.Errorf(
				"patch 'AreaID' cannot be an empty UUID: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if p.Notes.Set && p.Notes.Value != nil && *p.Notes.Value != nil {
		notesLength := len([]rune(strings.TrimSpace(**p.Notes.Value)))
		if notesLength > 2000 {
			return fmt.Errorf(
				"patch 'Notes' cannot exceed 2000 characters: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if p.Deadline.Set && p.Deadline.Value != nil && *p.Deadline.Value != nil {
		if (*p.Deadline.Value).IsZero() {
			return fmt.Errorf(
				"patch 'Deadline' cannot be a zero time: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}
