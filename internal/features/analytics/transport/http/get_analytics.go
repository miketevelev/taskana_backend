package analytics_transport_http

import (
	"net/http"

	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
	analytics_service "github.com/miketevelev/taskana_backend/internal/features/analytics/service"
)

func (h *AnalyticsHTTPHandler) GetAnalytics(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	startDate, err := core_http_request.GetDateQueryParam(r, "start_date")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid start_date")
		return
	}

	endDate, err := core_http_request.GetDateQueryParam(r, "end_date")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid end_date")
		return
	}

	projectID, err := core_http_request.GetUUIDQueryParam(r, "project_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid project_id")
		return
	}

	result, err := h.analyticsService.GetAnalytics(
		ctx, userID, analytics_service.AnalyticsFilter{
			StartDate: startDate,
			EndDate:   endDate,
			ProjectID: projectID,
		},
	)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get analytics")
		return
	}

	responseHandler.JSONResponse(result, http.StatusOK)
}
