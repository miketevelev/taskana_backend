package analytics_transport_http

import (
	"net/http"
	"strconv"

	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
	analytics_service "github.com/miketevelev/taskana_backend/internal/features/analytics/service"
)

func (h *AnalyticsHTTPHandler) GetYearlyAnalytics(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	// NOTE: unlike day/week/month, "year" is a plain integer (e.g. "2026"),
	// not a date, so it doesn't go through core_http_request.GetDateQueryParam.
	// If the codebase already has a GetIntQueryParam helper alongside
	// GetDateQueryParam/GetUUIDQueryParam, swap this block to use it instead.
	var year *int
	if raw := r.URL.Query().Get("year"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			responseHandler.ErrorResponse(err, "invalid year")
			return
		}
		year = &parsed
	}

	projectID, err := core_http_request.GetUUIDQueryParam(r, "project_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid project_id")
		return
	}

	result, err := h.analyticsService.GetYearlyAnalytics(
		ctx, userID, analytics_service.YearlyAnalyticsFilter{
			Year:      year,
			ProjectID: projectID,
		},
	)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get yearly analytics")
		return
	}

	responseHandler.JSONResponse(result, http.StatusOK)
}
