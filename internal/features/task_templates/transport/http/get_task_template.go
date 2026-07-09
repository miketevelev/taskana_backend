package task_templates_transport_http

import (
	"net/http"

	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
)

func (h *TaskTemplatesHTTPHandler) GetTaskTemplate(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	taskTemplateID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode task template request",
		)
		return
	}

	taskTemplate, err := h.taskTemplatesService.GetTaskTemplate(
		ctx, userID, taskTemplateID,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get task template",
		)
		return
	}

	response := taskTemplateDTOFromDomain(taskTemplate)

	responseHandler.JSONResponse(response, http.StatusOK)
}
