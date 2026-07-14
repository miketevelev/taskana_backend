package analytics_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_http_middleware "github.com/miketevelev/taskana_backend/internal/core/transport/http/middleware"
	core_http_server "github.com/miketevelev/taskana_backend/internal/core/transport/http/server"
	analytics_service "github.com/miketevelev/taskana_backend/internal/features/analytics/service"
)

type AnalyticsHTTPHandler struct {
	analyticsService AnalyticsService
	authMW           func(http.Handler) http.Handler
}

type AnalyticsService interface {
	GetAnalytics(
		ctx context.Context,
		userID uuid.UUID,
		filter analytics_service.AnalyticsFilter,
	) (domain.AnalyticsResponse, error)

	GetDailyAnalytics(
		ctx context.Context,
		userID uuid.UUID,
		filter analytics_service.DailyAnalyticsFilter,
	) (domain.AnalyticsResponse, error)

	GetWeeklyAnalytics(
		ctx context.Context,
		userID uuid.UUID,
		filter analytics_service.WeeklyAnalyticsFilter,
	) (domain.AnalyticsResponse, error)

	GetMonthlyAnalytics(
		ctx context.Context,
		userID uuid.UUID,
		filter analytics_service.MonthlyAnalyticsFilter,
	) (domain.AnalyticsResponse, error)

	GetYearlyAnalytics(
		ctx context.Context,
		userID uuid.UUID,
		filter analytics_service.YearlyAnalyticsFilter,
	) (domain.AnalyticsResponse, error)
}

func NewAnalyticsHTTPHandler(
	analyticsService AnalyticsService,
	tokenManager *core_auth.TokenManager,
) *AnalyticsHTTPHandler {
	return &AnalyticsHTTPHandler{
		analyticsService: analyticsService,
		authMW:           core_http_middleware.Auth(tokenManager),
	}
}

func (h *AnalyticsHTTPHandler) Routes() []core_http_server.Route {
	auth := []core_http_middleware.Middleware{
		func(next http.Handler) http.Handler { return h.authMW(next) },
	}

	return []core_http_server.Route{
		{
			Method:     http.MethodGet,
			Path:       "/analytics",
			Handler:    h.GetAnalytics,
			Middleware: auth,
		},
		{
			Method:     http.MethodGet,
			Path:       "/analytics/daily",
			Handler:    h.GetDailyAnalytics,
			Middleware: auth,
		},
		{
			Method:     http.MethodGet,
			Path:       "/analytics/weekly",
			Handler:    h.GetWeeklyAnalytics,
			Middleware: auth,
		},
		{
			Method:     http.MethodGet,
			Path:       "/analytics/monthly",
			Handler:    h.GetMonthlyAnalytics,
			Middleware: auth,
		},
		{
			Method:     http.MethodGet,
			Path:       "/analytics/yearly",
			Handler:    h.GetYearlyAnalytics,
			Middleware: auth,
		},
	}
}
