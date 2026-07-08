package task_templates_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_http_middleware "github.com/miketevelev/taskana_backend/internal/core/transport/http/middleware"
	core_http_server "github.com/miketevelev/taskana_backend/internal/core/transport/http/server"
)

type TaskTemplatesHTTPHandler struct {
	taskTemplatesService TaskTemplatesService
	authMW               func(http.Handler) http.Handler
}

type TaskTemplatesService interface {
	GetTaskTemplate(
		ctx context.Context,
		userID uuid.UUID,
		taskTemplateID uuid.UUID,
	) (domain.TaskTemplate, error)

	GetTaskTemplates(
		ctx context.Context,
		userID uuid.UUID,
		limit *int,
		offset *int,
	) ([]domain.TaskTemplate, error)

	CreateTaskTemplate(
		ctx context.Context,
		userID uuid.UUID,
		taskTemplate domain.TaskTemplate,
	) (domain.TaskTemplate, error)

	PatchTaskTemplate(
		ctx context.Context,
		userID uuid.UUID,
		taskTemplateID uuid.UUID,
		patch domain.TaskTemplatePatch,
	) (domain.TaskTemplate, error)
}

func NewTaskTemplatesHTTPHandler(
	taskTemplatesService TaskTemplatesService,
	tokenManager *core_auth.TokenManager,
) TaskTemplatesHTTPHandler {
	return TaskTemplatesHTTPHandler{
		taskTemplatesService: taskTemplatesService,
		authMW:               core_http_middleware.Auth(tokenManager),
	}
}

func (h *TaskTemplatesHTTPHandler) Routes() []core_http_server.Route {
	auth := []core_http_middleware.Middleware{
		func(next http.Handler) http.Handler { return h.authMW(next) },
	}

	return []core_http_server.Route{
		{
			Method:     http.MethodGet,
			Path:       "/task_templates/{id}",
			Handler:    h.GetTaskTemplate,
			Middleware: auth,
		},
		{
			Method:     http.MethodGet,
			Path:       "/task_templates",
			Handler:    h.GetTaskTemplates,
			Middleware: auth,
		},
		{
			Method:     http.MethodPost,
			Path:       "/task_templates",
			Handler:    h.CreateTaskTemplate,
			Middleware: auth,
		},
		{
			Method:     http.MethodPatch,
			Path:       "/task_templates/{id}",
			Handler:    h.PatchTaskTemplate,
			Middleware: auth,
		},
	}
}
