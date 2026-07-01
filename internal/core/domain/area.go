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
			"title must be between 1 and 255 characters long: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if a.Position < 0 {
		return fmt.Errorf(
			"position cannot be negative: %w", core_errors.ErrInvalidArgument,
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
