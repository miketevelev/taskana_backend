package timetracking_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_http_middleware "github.com/miketevelev/taskana_backend/internal/core/transport/http/middleware"
	core_http_server "github.com/miketevelev/taskana_backend/internal/core/transport/http/server"
)

type TimeTrackingHTTPHandler struct {
	timeTrackingService TimeTrackingService
	authMW              func(http.Handler) http.Handler
}

type TimeTrackingService interface {
	Sync(
		ctx context.Context,
		userID uuid.UUID,
		session []domain.PomodoroSession,
	) error
}

func NewTimeTrackingHTTPHandler(
	timeTrackingService TimeTrackingService,
	tokenManager *core_auth.TokenManager,
) *TimeTrackingHTTPHandler {
	return &TimeTrackingHTTPHandler{
		timeTrackingService: timeTrackingService,
		authMW:              core_http_middleware.Auth(tokenManager),
	}
}

func (h *TimeTrackingHTTPHandler) Routes() []core_http_server.Route {
	auth := []core_http_middleware.Middleware{
		func(next http.Handler) http.Handler { return h.authMW(next) },
	}

	return []core_http_server.Route{
		{
			Method:     http.MethodPost,
			Path:       "/timetracking/sync",
			Handler:    h.Sync,
			Middleware: auth,
		},
	}
}
