package checklists_transport_http

import (
	"fmt"
	"net/http"

	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
)

type ChangePositionChecklistRequest struct {
	Position int `json:"position" example:"2"`
}

func (r *ChangePositionChecklistRequest) Validate() error {
	if r.Position < 1 {
		return fmt.Errorf("'Position' must be 1 or greater")
	}
	return nil
}

func (h *ChecklistsHTTPHandler) ChangePosition(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	checklistID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to decode checklist request")
		return
	}

	var request ChangePositionChecklistRequest
	if err := core_http_request.DecodeAndValidateRequest(
		r, &request,
	); err != nil {
		responseHandler.ErrorResponse(
			err, "failed to decode change position request",
		)
		return
	}

	checklist, err := h.checklistsService.ChangePosition(
		ctx, userID, checklistID, request.Position,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to change checklist position",
		)
		return
	}

	response := checklistDTOFromDomain(checklist)

	responseHandler.JSONResponse(response, http.StatusOK)
}
