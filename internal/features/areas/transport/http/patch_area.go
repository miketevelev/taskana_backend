package areas_transport_http

import (
	"fmt"
	"net/http"

	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
	core_http_types "github.com/miketevelev/taskana_backend/internal/core/transport/http/types"
)

type PatchAreaRequest struct {
	Title    core_http_types.Nullable[string] `json:"title" example:"Home"`
	Position core_http_types.Nullable[int]    `json:"position" example:"2"`
}

func (r *PatchAreaRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("'Title' can't be NULL")
		}
		titleLength := len([]rune(*r.Title.Value))
		if titleLength < 3 || titleLength > 100 {
			return fmt.Errorf(
				"'Title' length must be between 3 and 100, got %d",
				titleLength,
			)
		}
	}

	if r.Position.Set {
		if r.Position.Value == nil {
			return fmt.Errorf("'Position' can't be NULL")
		}
		if *r.Position.Value < 1 {
			return fmt.Errorf("'Position' can't be negative")
		}
	}

	return nil
}

type PatchAreaResponse AreaDTOResponse

func (h *AreasHTTPHandler) PatchArea(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	areaID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode area request",
		)
		return
	}

	var request PatchAreaRequest
	if err := core_http_request.DecodeAndValidateRequest(
		r, &request,
	); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode patch area request",
		)
		return
	}

	areaPatch := areaPatchFromRequest(request)

	areaDomain, err := h.areasService.PatchArea(ctx, userID, areaID, areaPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch area",
		)
		return
	}

	response := PatchAreaResponse(areaDTOFromDomain(areaDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func areaPatchFromRequest(request PatchAreaRequest) domain.AreaPatch {
	return domain.NewAreaPatch(
		request.Title.ToDomain(),
		request.Position.ToDomain(),
	)
}
