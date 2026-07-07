package task_templates_transport_http

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
	task_templates_service "github.com/miketevelev/taskana_backend/internal/features/task_templates/service"
)

type CreateTaskTemplateRequest struct {
	ProjectID          *uuid.UUID `json:"project_id"`
	HeadingID          *uuid.UUID `json:"heading_id"`
	Title              string     `json:"title" validate:"required,min=3,max=100"`
	Notes              *string    `json:"notes"`
	RecurrenceRule     string     `json:"recurrence_rule" validate:"required"`
	RecurrenceType     string     `json:"recurrence_type" validate:"required"`
	TargetBucket       string     `json:"target_bucket" validate:"required"`
	NextExecutionDate  time.Time  `json:"next_execution_date" validate:"required"`
	IsTimeTracked      bool       `json:"is_time_tracked"`
	EstimatedPomodoros int        `json:"estimated_pomodoros"`
}

type CreateTaskTemplateResponse TaskTemplateDTOResponse

func (h *TaskTemplatesHTTPHandler) CreateTaskTemplate(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	var request CreateTaskTemplateRequest
	if err := core_http_request.DecodeAndValidateRequest(
		r, &request,
	); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode task template request",
		)
		return
	}

	recurrenceType, err := task_templates_service.
		ParseRecurrenceType(request.RecurrenceType)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"invalid recurrence_type",
		)
		return
	}

	targetBucket, err := task_templates_service.
		ParseTargetBucket(request.TargetBucket)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"invalid target_bucket",
		)
		return
	}

	taskTemplateDomain := domainFromDTO(
		userID,
		request,
		recurrenceType,
		targetBucket,
	)

	taskTemplate, err := h.taskTemplatesService.CreateTaskTemplate(
		ctx,
		userID, taskTemplateDomain,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create task template",
		)
		return
	}

	response := CreateTaskTemplateResponse(taskTemplateDTOFromDomain(taskTemplate))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(
	userID uuid.UUID,
	dto CreateTaskTemplateRequest,
	recurrenceType domain.RecurrenceType,
	targetBucket domain.TargetBucket,
) domain.TaskTemplate {
	return domain.NewTaskTemplateUninitialized(
		userID,
		dto.ProjectID,
		dto.HeadingID,
		dto.Title,
		dto.Notes,
		dto.RecurrenceRule,
		recurrenceType,
		targetBucket,
		dto.NextExecutionDate,
		dto.IsTimeTracked,
		dto.EstimatedPomodoros,
	)
}
