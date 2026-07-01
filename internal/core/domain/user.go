package domain

import (
	"fmt"
	"regexp"
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
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	emailLength := len(u.Email)
	if emailLength < 5 || emailLength > 254 {
		return fmt.Errorf(
			"email must be between 5 and 254 characters long, got %d: %w",
			emailLength, core_errors.ErrInvalidArgument,
		)
	}
	if !emailRegex.MatchString(u.Email) {
		return fmt.Errorf(
			"invalid email format: %w", core_errors.ErrInvalidArgument,
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

func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate users patch: %w", err)
	}

	tmp := *u

	if patch.FirstName.Set {
		tmp.FirstName = *patch.FirstName.Value
	}

	if patch.LastName.Set {
		tmp.LastName = *patch.LastName.Value
	}

	if patch.Email.Set {
		tmp.Email = *patch.Email.Value
	}

	if patch.Timezone.Set {
		tmp.Timezone = *patch.Timezone.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate users patch: %w", err)
	}

	*u = tmp

	return nil
}

type UserPatch struct {
	FirstName Nullable[string]
	LastName  Nullable[string]
	Email     Nullable[string]
	Timezone  Nullable[string]
}

func NewUserPatch(
	firstName Nullable[string],
	lastName Nullable[string],
	email Nullable[string],
	timezone Nullable[string],
) UserPatch {
	return UserPatch{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Timezone:  timezone,
	}
}

func (p *UserPatch) Validate() error {
	if p.FirstName.Set && p.FirstName.Value == nil {
		return fmt.Errorf(
			"first name cannot be nil: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.LastName.Set && p.LastName.Value == nil {
		return fmt.Errorf(
			"last name cannot be nil: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.Email.Set && p.Email.Value == nil {
		return fmt.Errorf(
			"email cannot be nil: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.Timezone.Set && p.Timezone.Value == nil {
		return fmt.Errorf(
			"timezone cannot be nil: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}
