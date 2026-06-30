package auth_transport_http

import (
	"time"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type RefreshAndLogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id" example:"10"`
	Version   int       `json:"version" example:"1"`
	FirstName string    `json:"first_name" example:"John"`
	LastName  string    `json:"last_name" example:"Doe"`
	Email     string    `json:"email" example:"mail@mail.com"`
	Timezone  string    `json:"timezone" example:"Europe/London"`
	CreatedAt time.Time `json:"created_at" example:"2020-01-01T00:00:00+00:00"`
	UpdatedAt time.Time `json:"updated_at" example:"2020-01-01T00:00:00+00:00"`
}

type RegisterAndLoginResponse struct {
	Tokens domain.TokenPair `json:"tokens"`
	User   UserResponse     `json:"user"`
}
