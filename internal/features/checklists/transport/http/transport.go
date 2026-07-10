package checklists_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_http_middleware "github.com/miketevelev/taskana_backend/internal/core/transport/http/middleware"
	core_http_server "github.com/miketevelev/taskana_backend/internal/core/transport/http/server"
)

type ChecklistsHTTPHandler struct {
	checklistsService ChecklistsService
	authMW            func(http.Handler) http.Handler
}

type ChecklistsService interface {
	CreateChecklist(
		ctx context.Context,
		userID uuid.UUID,
		checklist domain.Checklist,
	) (domain.Checklist, error)
}

func NewChecklistsHTTPHandler(
	checklistsService ChecklistsService,
	tokenManager *core_auth.TokenManager,
) ChecklistsHTTPHandler {
	return ChecklistsHTTPHandler{
		checklistsService: checklistsService,
		authMW:            core_http_middleware.Auth(tokenManager),
	}
}

func (h *ChecklistsHTTPHandler) Routes() []core_http_server.Route {
	auth := []core_http_middleware.Middleware{
		func(next http.Handler) http.Handler { return h.authMW(next) },
	}

	return []core_http_server.Route{
		{
			Method:     http.MethodPost,
			Path:       "/checklists",
			Handler:    h.CreateChecklist,
			Middleware: auth,
		},
	}
}
