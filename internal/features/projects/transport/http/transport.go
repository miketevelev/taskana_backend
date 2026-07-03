package projects_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_http_middleware "github.com/miketevelev/taskana_backend/internal/core/transport/http/middleware"
	core_http_server "github.com/miketevelev/taskana_backend/internal/core/transport/http/server"
)

type ProjectsHTTPHandler struct {
	projectsService ProjectsService
	authMW          func(http.Handler) http.Handler
}

type ProjectsService interface {
	CreateProject(
		ctx context.Context,
		userID uuid.UUID,
		project domain.Project,
	) (domain.Project, error)
}

func NewProjectsHTTPHandler(
	projectsService ProjectsService,
	tokenManager *core_auth.TokenManager,
) ProjectsHTTPHandler {
	return ProjectsHTTPHandler{
		projectsService: projectsService,
		authMW:          core_http_middleware.Auth(tokenManager),
	}
}

func (h *ProjectsHTTPHandler) Routes() []core_http_server.Route {
	auth := []core_http_middleware.Middleware{
		func(next http.Handler) http.Handler { return h.authMW(next) },
	}

	return []core_http_server.Route{
		{
			Method:     http.MethodPost,
			Path:       "/projects",
			Handler:    h.CreateProject,
			Middleware: auth,
		},
	}
}
