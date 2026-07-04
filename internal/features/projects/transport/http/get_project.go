package projects_transport_http

import (
	"net/http"

	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
)

func (h *ProjectsHTTPHandler) GetProject(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	projectID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode project request",
		)
		return
	}

	project, err := h.projectsService.GetProject(ctx, userID, projectID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get project",
		)
		return
	}

	response := projectDTOFromDomain(project)

	responseHandler.JSONResponse(response, http.StatusOK)
}
