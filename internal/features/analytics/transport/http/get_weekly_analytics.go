package analytics_transport_http

import (
	"net/http"

	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
	analytics_service "github.com/miketevelev/taskana_backend/internal/features/analytics/service"
)

func (h *AnalyticsHTTPHandler) GetWeeklyAnalytics(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	// Any date within the desired week; omit for the current week.
	week, err := core_http_request.GetDateQueryParam(r, "week")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid week")
		return
	}

	projectID, err := core_http_request.GetUUIDQueryParam(r, "project_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid project_id")
		return
	}

	result, err := h.analyticsService.GetWeeklyAnalytics(
		ctx, userID, analytics_service.WeeklyAnalyticsFilter{
			Week:      week,
			ProjectID: projectID,
		},
	)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get weekly analytics")
		return
	}

	responseHandler.JSONResponse(result, http.StatusOK)
}
