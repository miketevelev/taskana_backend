package user_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_http_middleware "github.com/miketevelev/taskana_backend/internal/core/transport/http/middleware"
	core_http_server "github.com/miketevelev/taskana_backend/internal/core/transport/http/server"
)

type UsersHTTPHandler struct {
	userService UserService
	authMW      func(http.Handler) http.Handler
}

type UserService interface {
	GetUser(
		ctx context.Context,
		userID uuid.UUID,
	) (domain.User, error)

	PatchUser(
		ctx context.Context,
		userID uuid.UUID,
		patch domain.UserPatch,
	) (domain.User, error)

	//DeleteUser(
	//	ctx context.Context,
	//) error
}

func NewUsersHTTPHandler(
	userService UserService,
	tokenManager *core_auth.TokenManager,
) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		userService: userService,
		authMW:      core_http_middleware.Auth(tokenManager),
	}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	auth := []core_http_middleware.Middleware{
		func(next http.Handler) http.Handler { return h.authMW(next) },
	}

	return []core_http_server.Route{
		{
			Method:     http.MethodGet,
			Path:       "/user",
			Handler:    h.GetUser,
			Middleware: auth,
		},
		{
			Method:     http.MethodPatch,
			Path:       "/user",
			Handler:    h.PatchUser,
			Middleware: auth,
		},
	}
}
