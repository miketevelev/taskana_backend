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
	GetChecklist(
		ctx context.Context,
		userID uuid.UUID,
		checklistID uuid.UUID,
	) (domain.Checklist, error)

	GetChecklists(
		ctx context.Context,
		userID uuid.UUID,
		limit *int,
		offset *int,
	) ([]domain.Checklist, error)

	CreateChecklist(
		ctx context.Context,
		userID uuid.UUID,
		checklist domain.Checklist,
	) (domain.Checklist, error)

	ChangePosition(
		ctx context.Context,
		userID uuid.UUID,
		checklistID uuid.UUID,
		newPosition int,
	) (domain.Checklist, error)

	PatchChecklist(
		ctx context.Context,
		userID uuid.UUID,
		checklistID uuid.UUID,
		patch domain.ChecklistPatch,
	) (domain.Checklist, error)

	//DeleteChecklist(
	//	ctx context.Context,
	//	userID uuid.UUID,
	//	checklistID uuid.UUID,
	//) (domain.Checklist, error)
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
			Method:     http.MethodGet,
			Path:       "/checklists/{id}",
			Handler:    h.GetChecklist,
			Middleware: auth,
		},
		{
			Method:     http.MethodGet,
			Path:       "/checklists",
			Handler:    h.GetChecklists,
			Middleware: auth,
		},
		{
			Method:     http.MethodPost,
			Path:       "/checklists",
			Handler:    h.CreateChecklist,
			Middleware: auth,
		},
		{
			Method:     http.MethodPost,
			Path:       "/checklists/{id}",
			Handler:    h.ChangePosition,
			Middleware: auth,
		},
		{
			Method:     http.MethodPatch,
			Path:       "/checklists/{id}",
			Handler:    h.PatchChecklist,
			Middleware: auth,
		},
		//{
		//	Method:     http.MethodDelete,
		//	Path:       "/checklists/{id}",
		//	Handler:    h.DeleteChecklist,
		//	Middleware: auth,
		//},
	}
}
