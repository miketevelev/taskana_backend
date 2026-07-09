package tasks_transport_http

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
	core_http_types "github.com/miketevelev/taskana_backend/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	ProjectID          core_http_types.Nullable[*uuid.UUID]        `json:"project_id"`
	HeadingID          core_http_types.Nullable[*uuid.UUID]        `json:"heading_id"`
	Title              core_http_types.Nullable[string]            `json:"title"`
	Notes              core_http_types.Nullable[*string]           `json:"notes"`
	Status             core_http_types.Nullable[domain.TaskStatus] `json:"status"`
	Bucket             core_http_types.Nullable[domain.TaskBucket] `json:"bucket"`
	StartDate          core_http_types.Nullable[*time.Time]        `json:"start_date"`
	Deadline           core_http_types.Nullable[*time.Time]        `json:"deadline"`
	IsTimeTracked      core_http_types.Nullable[bool]              `json:"is_time_tracked"`
	EstimatedPomodoros core_http_types.Nullable[int]               `json:"estimated_pomodoros"`
	CompletedAt        core_http_types.Nullable[*time.Time]        `json:"completed_at"`
}

func (r *PatchTaskRequest) Validate() error {
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
		case domain.TaskStatusOpen,
			domain.TaskStatusCompleted,
			domain.TaskStatusCanceled:
		default:
			return fmt.Errorf("invalid task status '%s'", status)
		}
	}

	if r.Bucket.Set {
		if r.Bucket.Value == nil {
			return fmt.Errorf("'Bucket' can't be NULL")
		}
		bucket := *r.Bucket.Value
		switch bucket {
		case domain.TaskBucketInbox,
			domain.TaskBucketToday,
			domain.TaskBucketAnytime,
			domain.TaskBucketSomeday:
		default:
			return fmt.Errorf("invalid task bucket '%s'", bucket)
		}
	}

	if r.ProjectID.Set && r.ProjectID.Value != nil && *r.ProjectID.Value != nil {
		if **r.ProjectID.Value == uuid.Nil {
			return fmt.Errorf("'ProjectID' cannot be an empty UUID")
		}
	}

	if r.HeadingID.Set && r.HeadingID.Value != nil && *r.HeadingID.Value != nil {
		if **r.HeadingID.Value == uuid.Nil {
			return fmt.Errorf("'HeadingID' cannot be an empty UUID")
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

	if r.EstimatedPomodoros.Set && r.EstimatedPomodoros.Value != nil {
		if *r.EstimatedPomodoros.Value < 0 {
			return fmt.Errorf("'EstimatedPomodoros' cannot be negative")
		}
	}

	if r.StartDate.Set && r.StartDate.Value != nil && *r.StartDate.Value != nil {
		if (**r.StartDate.Value).IsZero() {
			return fmt.Errorf("'StartDate' cannot be a zero time")
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

type PatchTaskResponse TaskDTOResponse

func (h *TasksHTTPHandler) PatchTask(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	taskID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get task ID path value",
		)
		return
	}

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(
		r, &request,
	); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode task request",
		)
		return
	}

	taskPatch := taskPatchFromRequest(request)

	taskDomain, err := h.tasksService.PatchTask(
		ctx, userID, taskID, taskPatch,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch task",
		)
		return
	}

	response := PatchTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.ProjectID.ToDomain(),
		request.HeadingID.ToDomain(),
		request.Title.ToDomain(),
		request.Notes.ToDomain(),
		request.Status.ToDomain(),
		request.Bucket.ToDomain(),
		request.StartDate.ToDomain(),
		request.Deadline.ToDomain(),
		request.IsTimeTracked.ToDomain(),
		request.EstimatedPomodoros.ToDomain(),
		request.CompletedAt.ToDomain(),
	)
}
