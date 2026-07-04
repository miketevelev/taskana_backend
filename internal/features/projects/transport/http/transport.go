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
	GetProject(
		ctx context.Context,
		userID uuid.UUID,
		projectID uuid.UUID,
	) (domain.Project, error)

	GetProjects(
		ctx context.Context,
		userID uuid.UUID,
		limit *int,
		offset *int,
	) ([]domain.Project, error)

	CreateProject(
		ctx context.Context,
		userID uuid.UUID,
		project domain.Project,
	) (domain.Project, error)

	ChangePosition(
		ctx context.Context,
		userID uuid.UUID,
		projectID uuid.UUID,
		newPosition int,
	) (domain.Project, error)

	PatchProject(
		ctx context.Context,
		userID uuid.UUID,
		projectID uuid.UUID,
		patch domain.ProjectPatch,
	) (domain.Project, error)

	DeleteProject(
		ctx context.Context,
		userID uuid.UUID,
		projectID uuid.UUID,
	) error
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
			Method:     http.MethodGet,
			Path:       "/projects/{id}",
			Handler:    h.GetProject,
			Middleware: auth,
		},
		{
			Method:     http.MethodGet,
			Path:       "/projects",
			Handler:    h.GetProjects,
			Middleware: auth,
		},
		{
			Method:     http.MethodPost,
			Path:       "/projects",
			Handler:    h.CreateProject,
			Middleware: auth,
		},
		{
			Method:     http.MethodPost,
			Path:       "/projects/{id}",
			Handler:    h.ChangePosition,
			Middleware: auth,
		},
		{
			Method:     http.MethodPatch,
			Path:       "/projects/{id}",
			Handler:    h.PatchProject,
			Middleware: auth,
		},
		{
			Method:     http.MethodDelete,
			Path:       "/projects/{id}",
			Handler:    h.DeleteProject,
			Middleware: auth,
		},
	}
}
