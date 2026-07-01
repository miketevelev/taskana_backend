package areas_transport_http

import (
	"net/http"

	"github.com/google/uuid"
	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
)

type CreateAreaRequest struct {
	Title string `json:"title" validate:"required,min=3,max=100" example:"Home"`
}

type CreateAreaResponse AreaDTOResponse

func (h *AreasHTTPHandler) CreateArea(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	var request CreateAreaRequest
	if err := core_http_request.DecodeAndValidateRequest(
		r, &request,
	); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode area request",
		)
		return
	}

	areaDomain := domainFromDTO(userID, request)

	area, err := h.areasService.CreateArea(ctx, userID, areaDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create area",
		)
		return
	}

	response := CreateAreaResponse(areaDTOFromDomain(area))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(
	userID uuid.UUID,
	dto CreateAreaRequest,
) domain.Area {
	return domain.NewAreaUninitialized(
		userID,
		dto.Title,
	)
}
