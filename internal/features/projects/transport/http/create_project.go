package projects_transport_http

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
)

type CreateProjectRequest struct {
	AreaID   *uuid.UUID `json:"area_id,omitempty"`
	Title    string     `json:"title"`
	Notes    *string    `json:"notes,omitempty"`
	Deadline *time.Time `json:"deadline,omitempty"`
}

type CreateProjectResponse ProjectDTOResponse

func (h *ProjectsHTTPHandler) CreateProject(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	var request CreateProjectRequest
	if err := core_http_request.DecodeAndValidateRequest(
		r, &request,
	); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode project request",
		)
		return
	}

	projectDomain := domainFromDTO(userID, request)

	project, err := h.projectsService.CreateProject(
		ctx,
		userID,
		projectDomain,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create project",
		)
		return
	}

	response := CreateProjectResponse(projectDTOFromDomain(project))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(
	userID uuid.UUID,
	dto CreateProjectRequest,
) domain.Project {
	return domain.NewProjectUninitialized(
		userID,
		dto.AreaID,
		dto.Title,
		dto.Notes,
		dto.Deadline,
	)
}
