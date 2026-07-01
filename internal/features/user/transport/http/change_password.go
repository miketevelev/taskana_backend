package user_transport_http

import (
	"errors"
	"net/http"
	"strings"

	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
	"github.com/miketevelev/taskana_backend/internal/features/auth/transport/http"
)

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=8,max=72" example:"12345678"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=72" example:"12345678"`
}

func (h *UsersHTTPHandler) ChangePassword(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	var request ChangePasswordRequest
	if err := core_http_request.DecodeAndValidateRequest(
		r, &request,
	); err != nil {
		responseHandler.ErrorResponse(
			err,
			"invalid password request",
		)
		return
	}

	userAgent := r.Header.Get("User-Agent")
	var ua *string
	if userAgent != "" {
		ua = &userAgent
	}

	tokens, user, err := h.userService.ChangePassword(
		ctx,
		userID,
		request.OldPassword,
		request.NewPassword,
		ua,
	)
	if err != nil {
		if errors.Is(err, core_errors.ErrInvalidArgument) {
			errMsg := err.Error()
			suffix := ": " + core_errors.ErrInvalidArgument.Error()
			errMsg = strings.TrimSuffix(errMsg, suffix)

			responseHandler.JSONResponse(
				map[string]string{
					"error":   "invalid password request",
					"message": errMsg,
				},
				http.StatusBadRequest,
			)
			return
		}
		if errors.Is(err, core_errors.ErrUnauthorized) {
			responseHandler.JSONResponse(
				map[string]string{
					"error":   "invalid password request",
					"message": "incorrect password",
				},
				http.StatusUnauthorized,
			)
			return
		}
		responseHandler.ErrorResponse(
			err,
			"invalid password request",
		)
		return
	}

	response := auth_transport_http.RegisterAndLoginResponse{
		Tokens: tokens,
		User: auth_transport_http.UserDTOResponse{
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
