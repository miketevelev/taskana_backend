package task_templates_transport_http

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

type PatchTaskTemplateRequest struct {
	ProjectID          core_http_types.Nullable[*uuid.UUID]            `json:"project_id,omitempty"`
	HeadingID          core_http_types.Nullable[*uuid.UUID]            `json:"heading_id,omitempty"`
	Title              core_http_types.Nullable[string]                `json:"title"`
	Notes              core_http_types.Nullable[*string]               `json:"notes,omitempty"`
	RecurrenceRule     core_http_types.Nullable[string]                `json:"recurrence_rule"`
	RecurrenceType     core_http_types.Nullable[domain.RecurrenceType] `json:"recurrence_type"`
	TargetBucket       core_http_types.Nullable[domain.TargetBucket]   `json:"target_bucket"`
	NextExecutionDate  core_http_types.Nullable[time.Time]             `json:"next_execution_date"`
	IsTimeTracked      core_http_types.Nullable[bool]                  `json:"is_time_tracked"`
	EstimatedPomodoros core_http_types.Nullable[int]                   `json:"estimated_pomodoros"`
}

func (r *PatchTaskTemplateRequest) Validate() error {
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

	if r.Notes.Set && r.Notes.Value != nil && *r.Notes.Value != nil {
		noteLength := len([]rune(strings.TrimSpace(**r.Notes.Value)))
		if noteLength > 2000 {
			return fmt.Errorf("'Notes' cannot exceed 2000 characters, got %d", noteLength)
		}
	}

	if r.RecurrenceRule.Set {
		if r.RecurrenceRule.Value == nil || strings.TrimSpace(*r.RecurrenceRule.Value) == "" {
			return fmt.Errorf("'RecurrenceRule' can't be NULL or empty string")
		}
	}

	if r.RecurrenceType.Set {
		if r.RecurrenceType.Value == nil {
			return fmt.Errorf("'RecurrenceType' can't be NULL")
		}
		switch *r.RecurrenceType.Value {
		case domain.RecurrenceTypeFixed, domain.RecurrenceTypeFromCompletion:
		default:
			return fmt.Errorf("invalid 'RecurrenceType' '%s''", *r.RecurrenceType.Value)
		}
	}

	if r.TargetBucket.Set {
		if r.TargetBucket.Value == nil {
			return fmt.Errorf("'TargetBucket' can't be NULL")
		}
		switch *r.TargetBucket.Value {
		case domain.TargetBucketToday, domain.TargetBucketInbox:
		default:
			return fmt.Errorf("invalid 'TargetBucket' '%s''", *r.TargetBucket.Value)
		}
	}

	if r.NextExecutionDate.Set {
		if r.NextExecutionDate.Value == nil {
			return fmt.Errorf("'NextExecutionDate' can't be NULL")
		}
		if (*r.NextExecutionDate.Value).IsZero() {
			return fmt.Errorf("'NextExecutionDate' can't be zero")
		}
	}

	if r.IsTimeTracked.Set && r.IsTimeTracked.Value == nil {
		return fmt.Errorf("'IsTimeTracked' can't be NULL")
	}

	if r.EstimatedPomodoros.Set {
		if r.EstimatedPomodoros.Value == nil {
			return fmt.Errorf("'EstimatedPomodoros' can't be NULL")
		}
		if *r.EstimatedPomodoros.Value < 0 {
			return fmt.Errorf(
				"'EstimatedPomodoros' can't be negative, got %d", *r.EstimatedPomodoros.Value,
			)
		}
	}

	return nil
}

type PatchTaskTemplateResponse TaskTemplateDTOResponse

func (h *TaskTemplatesHTTPHandler) PatchTaskTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	taskTemplateID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get task template ID path value",
		)
		return
	}

	var request PatchTaskTemplateRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode task template request",
		)
		return
	}

	taskTemplatePatch := taskTemplatePatchFromRequest(request)

	taskTemplateDomain, err := h.taskTemplatesService.PatchTaskTemplate(
		ctx, userID, taskTemplateID, taskTemplatePatch,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch task template",
		)
		return
	}

	response := PatchTaskTemplateResponse(taskTemplateDTOFromDomain(taskTemplateDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func taskTemplatePatchFromRequest(request PatchTaskTemplateRequest) domain.TaskTemplatePatch {
	return domain.TaskTemplatePatch{
		ProjectID:          request.ProjectID.ToDomain(),
		HeadingID:          request.HeadingID.ToDomain(),
		Title:              request.Title.ToDomain(),
		Notes:              request.Notes.ToDomain(),
		RecurrenceRule:     request.RecurrenceRule.ToDomain(),
		RecurrenceType:     request.RecurrenceType.ToDomain(),
		TargetBucket:       request.TargetBucket.ToDomain(),
		NextExecutionDate:  request.NextExecutionDate.ToDomain(),
		IsTimeTracked:      request.IsTimeTracked.ToDomain(),
		EstimatedPomodoros: request.EstimatedPomodoros.ToDomain(),
	}
}
