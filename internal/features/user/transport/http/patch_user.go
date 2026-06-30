package user_transport_http

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_request "github.com/miketevelev/taskana_backend/internal/core/transport/http/request"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
	core_http_types "github.com/miketevelev/taskana_backend/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	FirstName core_http_types.Nullable[string] `json:"first_name" example:"John"`
	LastName  core_http_types.Nullable[string] `json:"last_name" example:"Doe"`
	Email     core_http_types.Nullable[string] `json:"email" example:"mail@mail.com"`
	Timezone  core_http_types.Nullable[string] `json:"timezone" example:"Europe/Moscow"`
}

func (r *PatchUserRequest) Validate() error {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if r.FirstName.Set {
		if r.FirstName.Value == nil {
			return fmt.Errorf("field 'first_name' can't be nil")
		}
		firstNameLength := len([]rune(*r.FirstName.Value))
		if firstNameLength < 3 || firstNameLength > 100 {
			return fmt.Errorf(
				"field 'first_name' must be between 3 and 100"+
					" characters long, got %d", firstNameLength,
			)
		}
	}

	if r.LastName.Set {
		if r.LastName.Value == nil {
			return fmt.Errorf("field 'last_name' can't be nil")
		}
		lastNameLength := len([]rune(*r.LastName.Value))
		if lastNameLength < 3 || lastNameLength > 100 {
			return fmt.Errorf(
				"field 'last_name' must be between 3 and 100"+
					" characters long, got %d", lastNameLength,
			)
		}
	}

	if r.Email.Set {
		if r.Email.Value == nil {
			return fmt.Errorf("field 'email' can't be nil")
		}

		email := *r.Email.Value
		emailLength := len(email)

		if emailLength < 5 || emailLength > 254 {
			return fmt.Errorf(
				"field 'email' must be between 5 and 254 characters long, got %d",
				emailLength,
			)
		}

		if !emailRegex.MatchString(email) {
			return fmt.Errorf("field 'email' has invalid format")
		}
	}

	if r.Timezone.Set {
		if r.Timezone.Value == nil {
			return fmt.Errorf("field 'timezone' can't be nil")
		}

		tz := *r.Timezone.Value
		if _, err := time.LoadLocation(tz); err != nil {
			return fmt.Errorf(
				"invalid or unknown timezone '%s': %w",
				tz,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

func (h *UsersHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.MustUserIDFromContext(ctx)

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(
		r, &request,
	); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate PatchUser request",
		)
		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.userService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)

		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FirstName.ToDomain(),
		request.LastName.ToDomain(),
		request.Email.ToDomain(),
		request.Timezone.ToDomain(),
	)
}
