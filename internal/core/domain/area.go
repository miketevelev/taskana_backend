package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
)

type Area struct {
	ID      uuid.UUID `json:"id"`
	Version int       `json:"version"`

	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewArea(
	id uuid.UUID,
	version int,
	userID uuid.UUID,
	title string,
	position int,
	createdAt time.Time,
	updatedAt time.Time,
) Area {
	return Area{
		ID:        id,
		Version:   version,
		UserID:    userID,
		Title:     title,
		Position:  position,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func NewAreaUninitialized(
	userID uuid.UUID,
	title string,
) Area {
	now := time.Now().UTC()
	return Area{
		ID:        UninitializedID,
		Version:   UninitializedVersion,
		UserID:    userID,
		Title:     title,
		Position:  0,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (a *Area) Validate() error {
	titleLength := len([]rune(strings.TrimSpace(a.Title)))
	if titleLength < 3 || titleLength > 100 {
		return fmt.Errorf(
			"title must be between 3 and 255 characters long: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if a.Position < 1 && a.Position != 0 {
		return fmt.Errorf(
			"position must be 1 or greater: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if a.CreatedAt.IsZero() || a.UpdatedAt.IsZero() {
		return fmt.Errorf(
			"timestamps cannot be zero: %w", core_errors.ErrInvalidArgument,
		)
	}
	if a.UpdatedAt.Before(a.CreatedAt) {
		return fmt.Errorf(
			"updated_at cannot be before created_at: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if a.Version < 0 {
		a.Version = 0
	}

	return nil
}

func (a *Area) ApplyPatch(patch AreaPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate task patch: %w", err)
	}

	tmp := *a

	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}

	if patch.Position.Set {
		tmp.Position = *patch.Position.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate area patch: %w", err)
	}

	*a = tmp

	return nil
}

type AreaPatch struct {
	Title    Nullable[string]
	Position Nullable[int]
}

func NewAreaPatch(
	title Nullable[string],
	position Nullable[int],
) AreaPatch {
	return AreaPatch{
		Title:    title,
		Position: position,
	}
}

func (p *AreaPatch) Validate() error {
	if p.Title.Set && p.Title.Value == nil {
		return fmt.Errorf(
			"'Title' can't be patched to NULL: %w",
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

	return nil
}
