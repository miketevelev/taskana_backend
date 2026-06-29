package auth_transport_http

import (
	"net/http"

	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
)

type RegisterRequest struct {
	FirstName string `json:"first_name" validate:"required,min=3,max=100" example:"John"`
	LastName  string `json:"last_name" validate:"required,min=3,max=100" example:"Doe"`
	Email     string `json:"email" validate:"required,email" example:"mail@mail.com"`
	Password  string `json:"password" validate:"required,min=8" example:"123456"`
	Timezone  string `json:"timezone" validate:"required" example:"Europe/London"`
}

type RegisterResponse struct {
	Tokens domain.TokenPair `json:"tokens"`
	User   UserResponse     `json:"user"`
}

func (h *AuthHTTPHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	var request RegisterRequest
	if err := core_http_request.DecodeAndValidateRequest(
		r, &request,
	); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode register request",
		)
		return
	}

	tokens, user, err := h.authService.Register(
		ctx,
		request.FirstName,
		request.LastName,
		request.Email,
		request.Password,
		request.Timezone,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to register user",
		)
		return
	}

	response := RegisterResponse{
		Tokens: tokens,
		User: UserResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Timezone:  user.Timezone,
		},
	}

	responseHandler.JSONResponse(response, http.StatusCreated)
}
