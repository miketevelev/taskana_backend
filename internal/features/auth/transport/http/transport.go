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
	cleanup     func()
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

	Login(
		ctx context.Context,
		email string,
		password string,
		userAgent *string,
	) (domain.TokenPair, domain.User, error)

	Refresh(
		ctx context.Context,
		refreshToken string,
		userAgent *string,
	) (domain.TokenPair, error)
}

func NewAuthHTTPHandler(authService AuthService) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authService: authService,
	}
}

func (h *AuthHTTPHandler) Routes() []core_http_server.Route {
	rateLimit, stop := core_http_middleware.AuthRateLimit(
		5, time.Minute, 5*time.Minute,
	)
	h.cleanup = stop

	return []core_http_server.Route{
		{
			Method:     http.MethodPost,
			Path:       "/auth/register",
			Handler:    h.Register,
			Middleware: []core_http_middleware.Middleware{rateLimit},
		},
		{
			Method:     http.MethodPost,
			Path:       "/auth/login",
			Handler:    h.Login,
			Middleware: []core_http_middleware.Middleware{rateLimit},
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/refresh",
			Handler: h.Refresh,
		},
	}
}

func (h *AuthHTTPHandler) Shutdown() {
	if h.cleanup != nil {
		h.cleanup()
	}
}
