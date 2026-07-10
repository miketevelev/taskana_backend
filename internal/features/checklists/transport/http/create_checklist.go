package checklists_transport_http

import (
	"net/http"

	"github.com/google/uuid"
	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
)

type CreateChecklistRequest struct {
	TaskID uuid.UUID `json:"task_id"`
	Title  string    `json:"title"`
}

type CreateChecklistResponse ChecklistDTOResponse

func (h *ChecklistsHTTPHandler) CreateChecklist(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	var request CreateChecklistRequest
	if err := core_http_request.DecodeAndValidateRequest(
		r, &request,
	); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode checklist request",
		)
		return
	}

	checklistDomain := domainFromDTO(userID, request)

	checklist, err := h.checklistsService.CreateChecklist(
		ctx, userID, checklistDomain,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create checklist",
		)
		return
	}

	response := CreateChecklistResponse(checklistDTOFromDomain(checklist))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(
	userID uuid.UUID,
	dto CreateChecklistRequest,
) domain.Checklist {
	return domain.NewChecklistUninitialized(
		userID,
		dto.TaskID,
		dto.Title,
	)
}
