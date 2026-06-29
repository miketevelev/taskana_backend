package domain

import (
	"time"

	"github.com/google/uuid"
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
