package areas_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_http_middleware "github.com/miketevelev/taskana_backend/internal/core/transport/http/middleware"
	core_http_server "github.com/miketevelev/taskana_backend/internal/core/transport/http/server"
)

type AreasHTTPHandler struct {
	areasService AreasService
	authMW       func(http.Handler) http.Handler
}

type AreasService interface {
	CreateArea(
		ctx context.Context,
		userID uuid.UUID,
		area domain.Area,
	) (domain.Area, error)
}

func NewAreasHTTPHandler(
	areasService AreasService,
	tokenManager *core_auth.TokenManager,
) AreasHTTPHandler {
	return AreasHTTPHandler{
		areasService: areasService,
		authMW:       core_http_middleware.Auth(tokenManager),
	}
}

func (h *AreasHTTPHandler) Routes() []core_http_server.Route {
	auth := []core_http_middleware.Middleware{
		func(next http.Handler) http.Handler { return h.authMW(next) },
	}

	return []core_http_server.Route{
		{
			Method:     http.MethodPost,
			Path:       "/areas",
			Handler:    h.CreateArea,
			Middleware: auth,
		},
	}
}
