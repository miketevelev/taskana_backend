package headings_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_http_middleware "github.com/miketevelev/taskana_backend/internal/core/transport/http/middleware"
	core_http_server "github.com/miketevelev/taskana_backend/internal/core/transport/http/server"
)

type HeadingHTTPHandler struct {
	headingService HeadingService
	authMW         func(http.Handler) http.Handler
}

type HeadingService interface {
	CreateHeading(
		ctx context.Context,
		userID uuid.UUID,
		heading domain.Heading,
	) (domain.Heading, error)
}

func NewHeadingHTTPHandler(
	headingService HeadingService,
	tokenManager *core_auth.TokenManager,
) HeadingHTTPHandler {
	return HeadingHTTPHandler{
		headingService: headingService,
		authMW:         core_http_middleware.Auth(tokenManager),
	}
}

func (h *HeadingHTTPHandler) Routes() []core_http_server.Route {
	auth := []core_http_middleware.Middleware{
		func(next http.Handler) http.Handler { return h.authMW(next) },
	}

	return []core_http_server.Route{
		{
			Method:     http.MethodPost,
			Path:       "/headings",
			Handler:    h.CreateHeading,
			Middleware: auth,
		},
	}
}
