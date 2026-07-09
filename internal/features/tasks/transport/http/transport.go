package tasks_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_http_middleware "github.com/miketevelev/taskana_backend/internal/core/transport/http/middleware"
	core_http_server "github.com/miketevelev/taskana_backend/internal/core/transport/http/server"
)

type TasksHTTPHandler struct {
	tasksService TasksService
	authMW       func(http.Handler) http.Handler
}

type TasksService interface {
	CreateTask(
		ctx context.Context,
		userID uuid.UUID,
		task domain.Task,
	) (domain.Task, error)
}

func NewTasksHTTPHandler(
	tasksService TasksService,
	tokenManager *core_auth.TokenManager,
) TasksHTTPHandler {
	return TasksHTTPHandler{
		tasksService: tasksService,
		authMW:       core_http_middleware.Auth(tokenManager),
	}
}

func (h *TasksHTTPHandler) Routes() []core_http_server.Route {
	auth := []core_http_middleware.Middleware{
		func(next http.Handler) http.Handler { return h.authMW(next) },
	}

	return []core_http_server.Route{
		{
			Method:     http.MethodPost,
			Path:       "/tasks",
			Handler:    h.CreateTask,
			Middleware: auth,
		},
	}
}
