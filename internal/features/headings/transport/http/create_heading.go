package headings_transport_http

import (
	"net/http"

	"github.com/google/uuid"
	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
)

type CreateHeadingRequest struct {
	ProjectID uuid.UUID `json:"project_id"`
	Title     string    `json:"title"`
}

type CreateHeadingResponse HeadingDTOResponse

func (h *HeadingHTTPHandler) CreateHeading(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	var request CreateHeadingRequest
	if err := core_http_request.DecodeAndValidateRequest(
		r, &request,
	); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode heading request",
		)
		return
	}

	headingDomain := domainFromDTO(userID, request)

	heading, err := h.headingService.CreateHeading(
		ctx,
		userID,
		headingDomain,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create heading",
		)
		return
	}

	response := CreateHeadingResponse(headingDTOFromDomain(heading))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(
	userID uuid.UUID,
	dto CreateHeadingRequest,
) domain.Heading {
	return domain.NewHeadingUninitialized(
		userID,
		dto.ProjectID,
		dto.Title,
	)
}
