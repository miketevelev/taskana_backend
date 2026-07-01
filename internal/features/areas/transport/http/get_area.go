package areas_transport_http

import (
	"net/http"

	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
)

func (h *AreasHTTPHandler) GetArea(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	areaID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode area request",
		)
		return
	}

	area, err := h.areasService.GetArea(ctx, userID, areaID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get area",
		)
		return
	}

	response := areaDTOFromDomain(area)

	responseHandler.JSONResponse(response, http.StatusOK)
}
