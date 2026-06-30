package auth_transport_http

import (
	"net/http"

	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshResponse struct {
	Tokens domain.TokenPair `json:"tokens"`
}

func (h *AuthHTTPHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	var request RefreshRequest
	if err := core_http_request.DecodeAndValidateRequest(
		r, &request,
	); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode refresh request",
		)
	}

	userAgent := r.Header.Get("User-Agent")
	var ua *string
	if userAgent != "" {
		ua = &userAgent
	}

	tokens, err := h.authService.Refresh(ctx, request.RefreshToken, ua)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to refresh tokens",
		)
		return
	}

	response := RefreshResponse{
		Tokens: tokens,
	}

	responseHandler.JSONResponse(response, http.StatusOK)
}
