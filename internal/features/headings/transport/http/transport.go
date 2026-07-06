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
	GetHeading(
		ctx context.Context,
		userID uuid.UUID,
		headingID uuid.UUID,
	) (domain.Heading, error)

	GetHeadings(
		ctx context.Context,
		userID uuid.UUID,
		limit *int,
		offset *int,
	) ([]domain.Heading, error)

	CreateHeading(
		ctx context.Context,
		userID uuid.UUID,
		heading domain.Heading,
	) (domain.Heading, error)

	ChangePosition(
		ctx context.Context,
		userID uuid.UUID,
		headingID uuid.UUID,
		newPosition int,
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
			Method:     http.MethodGet,
			Path:       "/headings/{id}",
			Handler:    h.GetHeading,
			Middleware: auth,
		},
		{
			Method:     http.MethodGet,
			Path:       "/headings",
			Handler:    h.GetHeadings,
			Middleware: auth,
		},
		{
			Method:     http.MethodPost,
			Path:       "/headings",
			Handler:    h.CreateHeading,
			Middleware: auth,
		},
		{
			Method:     http.MethodPost,
			Path:       "/headings/{id}",
			Handler:    h.ChangePosition,
			Middleware: auth,
		},
	}
}
