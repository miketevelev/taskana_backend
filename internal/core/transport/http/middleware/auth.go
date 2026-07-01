package core_http_middleware

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
)

type TokenParser interface {
	ParseAccessToken(token string) (uuid.UUID, error)
}

func Auth(tokenManager *core_auth.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				log := core_logger.FromContext(ctx)
				responseHandler := core_http_response.NewHTTPResponseHandler(
					log, w,
				)

				authHeader := r.Header.Get("Authorization")
				if authHeader == "" {
					responseHandler.ErrorResponse(
						core_errors.ErrUnauthorized,
						"missing authorization header",
					)
					return
				}

				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
					responseHandler.ErrorResponse(
						core_errors.ErrUnauthorized,
						"invalid authorization header",
					)
					return
				}

				userID, err := tokenManager.ParseAccessToken(parts[1])
				if err != nil {
					responseHandler.ErrorResponse(err, "invalid access token")
					return
				}

				ctx = core_auth.WithUserID(ctx, userID)
				next.ServeHTTP(w, r.WithContext(ctx))
			},
		)
	}
}
