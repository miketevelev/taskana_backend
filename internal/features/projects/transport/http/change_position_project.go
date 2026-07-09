package projects_transport_http

import (
	"fmt"
	"net/http"

	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
)

type ChangePositionProjectRequest struct {
	Position int `json:"position" example:"2"`
}

func (r *ChangePositionProjectRequest) Validate() error {
	if r.Position < 1 {
		return fmt.Errorf("'Position' must be 1 or greater")
	}
	return nil
}

func (h *ProjectsHTTPHandler) ChangePosition(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	projectID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to decode project request")
		return
	}

	var request ChangePositionProjectRequest
	if err := core_http_request.DecodeAndValidateRequest(
		r, &request,
	); err != nil {
		responseHandler.ErrorResponse(
			err, "failed to decode change position request",
		)
		return
	}

	project, err := h.projectsService.ChangePosition(
		ctx, userID, projectID, request.Position,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to change project position",
		)
		return
	}

	response := projectDTOFromDomain(project)

	responseHandler.JSONResponse(response, http.StatusOK)
}
