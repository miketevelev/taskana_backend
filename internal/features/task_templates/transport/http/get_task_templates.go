package task_templates_transport_http

import (
	"net/http"

	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
)

type GetTaskTemplatesResponse []TaskTemplateDTOResponse

func (h *TaskTemplatesHTTPHandler) GetTaskTemplates(
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

	taskTemplates, err := h.taskTemplatesService.GetTaskTemplates(
		ctx, userID, limit, offset,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to fetch task templates",
		)
		return
	}

	response := GetTaskTemplatesResponse(taskTemplateDTOsFromDomains(taskTemplates))

	responseHandler.JSONResponse(response, http.StatusOK)
}
