package user_postgres_repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type UserModel struct {
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

func userDomainFromModel(userModel UserModel) domain.User {
	return domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FirstName,
		userModel.LastName,
		userModel.Email,
		userModel.PasswordHash,
		userModel.Timezone,
		userModel.CreatedAt,
		userModel.UpdatedAt,
	)
}
