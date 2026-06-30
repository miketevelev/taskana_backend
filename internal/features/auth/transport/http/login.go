package auth_transport_http

import (
	"net/http"
	"strings"

	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"mail@mail.com"`
	Password string `json:"password" validate:"required,min=8" example:"123456"`
}

func (h *AuthHTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	var request LoginRequest
	if err := core_http_request.DecodeAndValidateRequest(
		r, &request,
	); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode login request",
		)
		return
	}

	userAgent := r.Header.Get("User-Agent")
	var ua *string
	if userAgent != "" {
		ua = &userAgent
	}

	request.Email = strings.ToLower(strings.TrimSpace(request.Email))

	tokens, user, err := h.authService.Login(
		ctx,
		request.Email,
		request.Password,
		ua,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to login",
		)
		return
	}

	response := RegisterAndLoginResponse{
		Tokens: tokens,
		User: UserResponse{
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

	responseHandler.JSONResponse(response, http.StatusOK)
}
