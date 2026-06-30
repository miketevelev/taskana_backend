package auth_transport_http

import (
	"errors"
	"net/http"
	"strings"

	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
	"go.uber.org/zap"
)

type RegisterRequest struct {
	FirstName string `json:"first_name" validate:"required,min=3,max=100" example:"John"`
	LastName  string `json:"last_name" validate:"required,min=3,max=100" example:"Doe"`
	Email     string `json:"email" validate:"required,email" example:"mail@mail.com"`
	Password  string `json:"password" validate:"required,min=8" example:"123456"`
	Timezone  string `json:"timezone" validate:"required" example:"Europe/London"`
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

	request.Email = strings.ToLower(strings.TrimSpace(request.Email))

	tokens, user, err := h.authService.Register(
		ctx,
		request.FirstName,
		request.LastName,
		request.Email,
		request.Password,
		request.Timezone,
	)
	if err != nil {
		if errors.Is(err, core_errors.ErrAlreadyExists) {
			log.Error("failed to register user (duplicate)", zap.Error(err))

			responseHandler.JSONResponse(
				map[string]string{
					"error":   "failed to register user",
					"message": "email is already exists",
				},
				http.StatusConflict,
			)
			return
		}
		responseHandler.ErrorResponse(
			err,
			"failed to register user",
		)
		return
	}
	//if err != nil {
	//	responseHandler.ErrorResponse(
	//		err,
	//		"failed to register user",
	//	)
	//	return
	//}

	response := RegisterAndLoginResponse{
		Tokens: tokens,
		User: UserDTOResponse{
			ID:        user.ID,
			Version:   user.Version,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Timezone:  user.Timezone,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	responseHandler.JSONResponse(response, http.StatusCreated)
}
