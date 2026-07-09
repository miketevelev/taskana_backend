package headings_transport_http

import (
	"fmt"
	"net/http"

	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
)

type ChangePositionHeadingRequest struct {
	Position int `json:"position" example:"2"`
}

func (r *ChangePositionHeadingRequest) Validate() error {
	if r.Position < 1 {
		return fmt.Errorf("'Position' must be 1 or greater")
	}
	return nil
}

func (h *HeadingHTTPHandler) ChangePosition(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	headingID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to decode heading request")
		return
	}

	var request ChangePositionHeadingRequest
	if err := core_http_request.DecodeAndValidateRequest(
		r, &request,
	); err != nil {
		responseHandler.ErrorResponse(
			err, "failed to decode change position request",
		)
		return
	}

	heading, err := h.headingService.ChangePosition(
		ctx, userID, headingID, request.Position,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to change heading position",
		)
		return
	}

	response := headingDTOFromDomain(heading)

	responseHandler.JSONResponse(response, http.StatusOK)
}
