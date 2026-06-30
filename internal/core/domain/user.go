package domain

import (
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
)

type User struct {
	ID      uuid.UUID `json:"id"`
	Version int       `json:"version"`

	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Timezone     string    `json:"timezone"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewUser(
	id uuid.UUID,
	version int,
	firstName string,
	lastName string,
	email string,
	passwordHash string,
	timezone string,
	createdAt time.Time,
	updatedAt time.Time,
) User {
	return User{
		ID:           id,
		Version:      version,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: passwordHash,
		Timezone:     timezone,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}

func (u User) Validate() error {
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return fmt.Errorf(
			"invalid email format: %w", core_errors.ErrInvalidArgument,
		)
	}
	if len(u.Email) > 255 {
		return fmt.Errorf(
			"email is too long: %w", core_errors.ErrInvalidArgument,
		)
	}

	if strings.TrimSpace(u.FirstName) == "" || len(u.FirstName) > 100 {
		return fmt.Errorf(
			"invalid first name length: %w", core_errors.ErrInvalidArgument,
		)
	}
	if strings.TrimSpace(u.LastName) == "" || len(u.LastName) > 100 {
		return fmt.Errorf(
			"invalid last name length: %w", core_errors.ErrInvalidArgument,
		)
	}

	if len(u.PasswordHash) < 8 {
		return fmt.Errorf(
			"password must be at least 8 characters long: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if _, err := time.LoadLocation(u.Timezone); err != nil {
		return fmt.Errorf(
			"invalid or unknown timezone: %w", core_errors.ErrInvalidArgument,
		)
	}

	if u.CreatedAt.IsZero() || u.UpdatedAt.IsZero() {
		return fmt.Errorf(
			"timestamps cannot be zero: %w", core_errors.ErrInvalidArgument,
		)
	}
	if u.UpdatedAt.Before(u.CreatedAt) {
		return fmt.Errorf(
			"updated_at cannot be before created_at: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if u.Version < 0 {
		u.Version = 0
	}

	return nil
}

func NewUserUninitialized(
	firstName string,
	lastName string,
	email string,
	password string,
	timezone string,
) User {
	return NewUser(
		UninitializedID,
		UninitializedVersion,
		firstName,
		lastName,
		email,
		password,
		timezone,
		time.Now().UTC(),
		time.Now().UTC(),
	)
}
