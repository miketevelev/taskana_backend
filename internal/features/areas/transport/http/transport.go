package areas_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain/area"
	core_http_middleware "github.com/miketevelev/taskana_backend/internal/core/transport/http/middleware"
	core_http_server "github.com/miketevelev/taskana_backend/internal/core/transport/http/server"
)

type AreasHTTPHandler struct {
	areasService AreasService
	authMW       func(http.Handler) http.Handler
}

type AreasService interface {
	GetArea(
		ctx context.Context,
		userID uuid.UUID,
		areaID uuid.UUID,
	) (domain_area.Area, error)

	GetAreas(
		ctx context.Context,
		userID uuid.UUID,
		limit *int,
		offset *int,
	) ([]domain_area.Area, error)

	CreateArea(
		ctx context.Context,
		userID uuid.UUID,
		area domain_area.Area,
	) (domain_area.Area, error)

	ChangePosition(
		ctx context.Context,
		userID uuid.UUID,
		areaID uuid.UUID,
		newPosition int,
	) (domain_area.Area, error)

	PatchArea(
		ctx context.Context,
		userID uuid.UUID,
		areaID uuid.UUID,
		patch domain_area.AreaPatch,
	) (domain_area.Area, error)

	DeleteArea(
		ctx context.Context,
		userID uuid.UUID,
		areaID uuid.UUID,
	) error
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
			Method:     http.MethodGet,
			Path:       "/areas/{id}",
			Handler:    h.GetArea,
			Middleware: auth,
		},
		{
			Method:     http.MethodGet,
			Path:       "/areas",
			Handler:    h.GetAreas,
			Middleware: auth,
		},
		{
			Method:     http.MethodPost,
			Path:       "/areas",
			Handler:    h.CreateArea,
			Middleware: auth,
		},
		{
			Method:     http.MethodPost,
			Path:       "/areas/{id}",
			Handler:    h.ChangePosition,
			Middleware: auth,
		},
		{
			Method:     http.MethodPatch,
			Path:       "/areas/{id}",
			Handler:    h.PatchArea,
			Middleware: auth,
		},
		{
			Method:     http.MethodDelete,
			Path:       "/areas/{id}",
			Handler:    h.DeleteArea,
			Middleware: auth,
		},
	}
}
