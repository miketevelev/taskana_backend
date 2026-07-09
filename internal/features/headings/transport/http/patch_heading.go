package headings_transport_http

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

type PatchHeadingRequest struct {
	Title core_http_types.Nullable[string] `json:"title"`
}

func (r *PatchHeadingRequest) Validate() error {
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

	return nil
}

type PatchHeadingResponse HeadingDTOResponse

func (h *HeadingHTTPHandler) PatchHeading(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	headingID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get heading ID path value",
		)
		return
	}

	var request PatchHeadingRequest
	if err := core_http_request.DecodeAndValidateRequest(
		r, &request,
	); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode heading request",
		)
		return
	}

	headingPatch := headingPatchFromRequest(request)

	headingDomain, err := h.headingService.PatchHeading(
		ctx, userID, headingID, headingPatch,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch heading",
		)
		return
	}

	response := PatchHeadingResponse(headingDTOFromDomain(headingDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func headingPatchFromRequest(request PatchHeadingRequest) domain.HeadingPatch {
	return domain.NewHeadingPatch(
		request.Title.ToDomain(),
	)
}
