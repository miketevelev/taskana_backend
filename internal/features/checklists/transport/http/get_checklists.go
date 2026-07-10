package checklists_transport_http

import (
	"net/http"

	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
)

type GetChecklistsResponse []ChecklistDTOResponse

func (h *ChecklistsHTTPHandler) GetChecklists(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	limit, offset, err := core_http_request.GetLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get limit and offset query params",
		)
		return
	}

	checklists, err := h.checklistsService.GetChecklists(
		ctx, userID, limit, offset,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get checklists",
		)
		return
	}

	response := GetChecklistsResponse(checklistDTOsFromDomains(checklists))

	responseHandler.JSONResponse(response, http.StatusOK)
}
