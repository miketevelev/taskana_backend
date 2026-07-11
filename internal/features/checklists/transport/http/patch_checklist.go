package checklists_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
	core_http_types "github.com/miketevelev/taskana_backend/internal/core/transport/http/types"
)

type PatchChecklistRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"`
	IsCompleted core_http_types.Nullable[bool]   `json:"is_completed"`
}

func (r *PatchChecklistRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("'Title' can't be NULL")
		}
		titleLength := len([]rune(strings.TrimSpace(*r.Title.Value)))
		if titleLength < 3 || titleLength > 255 {
			return fmt.Errorf(
				"'Title' length must be between 3 and 255 characters, got %d",
				titleLength,
			)
		}
	}

	if r.IsCompleted.Set {
		if r.IsCompleted.Value == nil {
			return fmt.Errorf("'IsCompleted' can't be NULL")
		}
	}

	return nil
}

type PatchChecklistResponse ChecklistDTOResponse

func (h *ChecklistsHTTPHandler) PatchChecklist(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	checklistID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get checklist ID path value",
		)
		return
	}

	var request PatchChecklistRequest
	if err := core_http_request.DecodeAndValidateRequest(
		r, &request,
	); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode checklist request",
		)
		return
	}

	checklistPatch := checklistPatchFromRequest(request)

	checklistDomain, err := h.checklistsService.PatchChecklist(
		ctx, userID, checklistID, checklistPatch,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch checklist",
		)
		return
	}

	response := PatchChecklistResponse(checklistDTOFromDomain(checklistDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func checklistPatchFromRequest(request PatchChecklistRequest) domain.ChecklistPatch {
	return domain.NewChecklistPatch(
		request.Title.ToDomain(),
		request.IsCompleted.ToDomain(),
	)
}
