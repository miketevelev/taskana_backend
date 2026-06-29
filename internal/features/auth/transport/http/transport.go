package auth_transport_http

import (
	"context"
	"net/http"
	"time"

	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_http_middleware "github.com/miketevelev/taskana_backend/internal/core/transport/http/middleware"
	core_http_server "github.com/miketevelev/taskana_backend/internal/core/transport/http/server"
)

type AuthHTTPHandler struct {
	authService AuthService
}

type AuthService interface {
	Register(
		ctx context.Context,
		firstName string,
		lastName string,
		email string,
		password string,
		timezone string,
	) (domain.TokenPair, domain.User, error)
}

func NewAuthHTTPHandler(
	authService AuthService,
) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authService: authService,
	}
}

func (h *AuthHTTPHandler) Routes() []core_http_server.Route {
	rateLimit := core_http_middleware.AuthRateLimit(5, time.Minute)

	return []core_http_server.Route{
		{
			Method:     http.MethodPost,
			Path:       "/auth/register",
			Handler:    h.Register,
			Middleware: []core_http_middleware.Middleware{rateLimit},
		},
	}
}
