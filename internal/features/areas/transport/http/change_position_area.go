package areas_transport_http

import (
	"fmt"
	"net/http"

	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
)

type ChangePositionRequest struct {
	Position int `json:"position" example:"2"`
}

func (r *ChangePositionRequest) Validate() error {
	if r.Position < 1 {
		return fmt.Errorf("'Position' must be 1 or greater")
	}
	return nil
}

func (h *AreasHTTPHandler) ChangePosition(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	areaID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to decode area request")
		return
	}

	var request ChangePositionRequest
	if err := core_http_request.DecodeAndValidateRequest(
		r, &request,
	); err != nil {
		responseHandler.ErrorResponse(
			err, "failed to decode change position request",
		)
		return
	}

	area, err := h.areasService.ChangePosition(
		ctx, userID, areaID, request.Position,
	)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to change area position")
		return
	}

	response := areaDTOFromDomain(area)

	responseHandler.JSONResponse(response, http.StatusOK)
}
