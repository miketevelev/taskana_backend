package timetracking_transport_http

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

type SessionRequest struct {
	ID              uuid.UUID `json:"id"`
	TaskID          uuid.UUID `json:"task_id"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	DurationSeconds int       `json:"duration_seconds"`
	IsInterrupted   bool      `json:"is_interrupted"`
}

type SyncRequest struct {
	Sessions []SessionRequest `json:"sessions"`
}

func (h *TimeTrackingHTTPHandler) Sync(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	var request SyncRequest
	if err := core_http_request.DecodeAndValidateRequest(
		r, &request,
	); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode session request",
		)
		return
	}

	sessions := make([]domain.PomodoroSession, len(request.Sessions))
	for i, s := range request.Sessions {
		session := domain.NewPomodoroSessionUninitialized(
			s.ID,
			userID,
			s.TaskID,
			s.StartTime,
			s.EndTime,
			s.DurationSeconds,
			s.IsInterrupted,
		)
		sessions[i] = session
	}

	if err := h.timeTrackingService.Sync(
		ctx,
		userID,
		sessions,
	); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to sync sessions",
		)
		return
	}

	responseHandler.NoContentResponse()
}
