package tasks_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	ProjectID          *uuid.UUID        `json:"project_id,omitempty"`
	HeadingID          *uuid.UUID        `json:"heading_id,omitempty"`
	TemplateID         *uuid.UUID        `json:"template_id,omitempty"`
	Title              string            `json:"title"`
	Notes              *string           `json:"notes,omitempty"`
	Bucket             domain.TaskBucket `json:"bucket"`
	StartDate          *time.Time        `json:"start_date,omitempty"`
	Deadline           *time.Time        `json:"deadline,omitempty"`
	IsTimeTracked      bool              `json:"is_time_tracked"`
	EstimatedPomodoros int               `json:"estimated_pomodoros"`
}

func (r *CreateTaskRequest) Validate() error {
	if r.Bucket == "" {
		r.Bucket = domain.TaskBucketInbox
	}

	switch r.Bucket {
	case domain.TaskBucketInbox, domain.TaskBucketToday, domain.TaskBucketAnytime, domain.TaskBucketSomeday:
	default:
		return fmt.Errorf(
			"invalid task bucket '%s': %w", r.Bucket,
			core_errors.ErrInvalidArgument,
		)
	}

	if r.HeadingID != nil && r.ProjectID == nil {
		return fmt.Errorf(
			"task cannot have a heading without a project: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if r.StartDate != nil && r.Deadline != nil {
		if r.Deadline.Before(*r.StartDate) {
			return fmt.Errorf(
				"deadline cannot be before start_date: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}

type CreateTaskResponse TaskDTOResponse

func (h *TasksHTTPHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	var request CreateTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(
		r, &request,
	); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode task request",
		)
		return
	}

	taskDomain := domainFromCreateDTO(userID, request)

	task, err := h.tasksService.CreateTask(ctx, userID, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create task",
		)
		return
	}

	response := CreateTaskResponse(taskDTOFromDomain(task))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromCreateDTO(
	userID uuid.UUID,
	dto CreateTaskRequest,
) domain.Task {
	return domain.NewTaskUninitialized(
		userID,
		dto.ProjectID,
		dto.HeadingID,
		dto.TemplateID,
		dto.Title,
		dto.Notes,
		dto.Bucket,
		dto.StartDate,
		dto.Deadline,
		dto.IsTimeTracked,
		dto.EstimatedPomodoros,
	)
}
