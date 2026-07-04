package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
)

type Heading struct {
	ID      uuid.UUID `json:"id"`
	Version int       `json:"version"`

	UserID    uuid.UUID `json:"user_id"`
	ProjectID uuid.UUID `json:"project_id"`
	Title     string    `json:"title"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewHeading(
	id uuid.UUID,
	version int,
	userID uuid.UUID,
	projectID uuid.UUID,
	title string,
	position int,
	createdAt time.Time,
	updatedAt time.Time,
) Heading {
	return Heading{
		ID:        id,
		Version:   version,
		UserID:    userID,
		ProjectID: projectID,
		Title:     title,
		Position:  position,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func NewHeadingUninitialized(
	userID uuid.UUID,
	projectID uuid.UUID,
	title string,
) Heading {
	now := time.Now().UTC()
	return Heading{
		ID:        UninitializedID,
		Version:   UninitializedVersion,
		UserID:    userID,
		ProjectID: projectID,
		Title:     title,
		Position:  1,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (h *Heading) Validate() error {
	titleLength := len([]rune(strings.TrimSpace(h.Title)))
	if titleLength < 3 || titleLength > 100 {
		return fmt.Errorf(
			"title must be between 3 and 100 characters long: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if h.Position < 1 {
		return fmt.Errorf(
			"position must be 1 or greater: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if h.UserID == uuid.Nil {
		return fmt.Errorf(
			"user_id cannot be an empty UUID: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if h.ProjectID == uuid.Nil {
		return fmt.Errorf(
			"project_id cannot be an empty UUID: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if h.CreatedAt.IsZero() || h.UpdatedAt.IsZero() {
		return fmt.Errorf(
			"timestamps cannot be zero: %w", core_errors.ErrInvalidArgument,
		)
	}
	if h.UpdatedAt.Before(h.CreatedAt) {
		return fmt.Errorf(
			"updated_at cannot be before created_at: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if h.Version < 0 {
		h.Version = 0
	}

	return nil
}

func (h *Heading) ApplyPatch(patch HeadingPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate heading patch: %w", err)
	}

	tmp := *h

	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate heading after patch: %w", err)
	}

	*h = tmp

	return nil
}

type HeadingPatch struct {
	Title Nullable[string]
}

func NewHeadingPatch(
	title Nullable[string],
) HeadingPatch {
	return HeadingPatch{
		Title: title,
	}
}

func (p *HeadingPatch) Validate() error {
	if p.Title.Set && p.Title.Value == nil {
		return fmt.Errorf(
			"'Title' can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}
