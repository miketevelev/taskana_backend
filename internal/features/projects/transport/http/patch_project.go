package projects_transport_http

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain/project"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
	core_http_types "github.com/miketevelev/taskana_backend/internal/core/transport/http/types"
)

type PatchProjectRequest struct {
	AreaID      core_http_types.Nullable[*uuid.UUID]                   `json:"area_id"`
	Title       core_http_types.Nullable[string]                       `json:"title"`
	Notes       core_http_types.Nullable[*string]                      `json:"notes"`
	Status      core_http_types.Nullable[domain_project.ProjectStatus] `json:"status"`
	Deadline    core_http_types.Nullable[*time.Time]                   `json:"deadline"`
	CompletedAt core_http_types.Nullable[*time.Time]                   `json:"completed_at"`
}

func (r *PatchProjectRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("'Title' can't be NULL")
		}
		titleLength := len([]rune(strings.TrimSpace(*r.Title.Value)))
		if titleLength < 3 || titleLength > 100 {
			return fmt.Errorf(
				"'Title' length must be between 3 and 100 characters, got %d",
				titleLength,
			)
		}
	}

	if r.Status.Set {
		if r.Status.Value == nil {
			return fmt.Errorf("'Status' can't be NULL")
		}
		status := *r.Status.Value
		switch status {
		case domain_project.ProjectStatusActive, domain_project.ProjectStatusCompleted, domain_project.ProjectStatusDropped:
			// статус валиден
		default:
			return fmt.Errorf("invalid project status '%s'", status)
		}
	}

	if r.AreaID.Set && r.AreaID.Value != nil && *r.AreaID.Value != nil {
		if **r.AreaID.Value == uuid.Nil {
			return fmt.Errorf("'AreaID' cannot be an empty UUID")
		}
	}

	if r.Notes.Set && r.Notes.Value != nil && *r.Notes.Value != nil {
		notesLength := len([]rune(strings.TrimSpace(**r.Notes.Value)))
		if notesLength > 2000 {
			return fmt.Errorf(
				"'Notes' cannot exceed 2000 characters, got %d", notesLength,
			)
		}
	}

	if r.Deadline.Set && r.Deadline.Value != nil && *r.Deadline.Value != nil {
		if (**r.Deadline.Value).IsZero() {
			return fmt.Errorf("'Deadline' cannot be a zero time")
		}
	}

	if r.CompletedAt.Set && r.CompletedAt.Value != nil && *r.CompletedAt.Value != nil {
		if (**r.CompletedAt.Value).IsZero() {
			return fmt.Errorf("'CompletedAt' cannot be a zero time")
		}
	}

	return nil
}

type PatchProjectResponse ProjectDTOResponse

func (h *ProjectsHTTPHandler) PatchProject(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	projectID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get project ID path value",
		)
		return
	}

	var request PatchProjectRequest
	if err := core_http_request.DecodeAndValidateRequest(
		r, &request,
	); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode project request",
		)
		return
	}

	projectPatch := projectPatchFromRequest(request)

	projectDomain, err := h.projectsService.PatchProject(
		ctx, userID, projectID, projectPatch,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch project",
		)
		return
	}

	response := PatchProjectResponse(projectDTOFromDomain(projectDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func projectPatchFromRequest(request PatchProjectRequest) domain_project.ProjectPatch {
	return domain_project.NewProjectPatch(
		request.AreaID.ToDomain(),
		request.Title.ToDomain(),
		request.Notes.ToDomain(),
		request.Status.ToDomain(),
		request.Deadline.ToDomain(),
		request.CompletedAt.ToDomain(),
	)
}
