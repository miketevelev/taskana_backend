package core_http_middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
	"go.uber.org/zap"
)

const (
	requestIDHeader = "X-Request-Id"
)

func CORS() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				allowedOrigins := map[string]struct{}{
					"http://localhost:8080": {},
				}

				origin := r.Header.Get("Origin")
				if _, ok := allowedOrigins[origin]; !ok {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set(
						"Access-Control-Allow-Methods",
						"GET, POST, PUT, DELETE, OPTIONS, PATCH",
					)
					w.Header().Set(
						"Access-Control-Allow-Headers",
						"Content-Type, Authorization",
					)
					w.Header().Set("Vary", "Origin")
				}

				if r.Method == http.MethodOptions {
					w.WriteHeader(http.StatusNoContent)
					return
				}

				next.ServeHTTP(w, r)
			},
		)
	}
}

func RequestId() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				requestId := r.Header.Get(requestIDHeader)
				if requestId == "" {
					requestId = uuid.NewString()
				}

				r.Header.Set(requestIDHeader, requestId)
				w.Header().Set(requestIDHeader, requestId)

				next.ServeHTTP(w, r)
			},
		)
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				RequestId := r.Header.Get(requestIDHeader)
				l := log.With(
					zap.String("request_id", RequestId),
					zap.String("url", r.URL.String()),
				)

				ctx := core_logger.ToContext(r.Context(), l)

				next.ServeHTTP(w, r.WithContext(ctx))
			},
		)
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				log := core_logger.FromContext(ctx)

				rw := core_http_response.NewResponseWriter(w)

				before := time.Now()

				log.Debug(
					">>> incoming HTTP request",
					zap.String("method", r.Method),
					zap.Time("time", before.UTC()),
				)

				next.ServeHTTP(rw, r)

				log.Debug(
					"<<< done HTTP request",
					zap.Int("status_code", rw.GetStatusCode()),
					zap.Duration("latency", time.Now().Sub(before)),
				)
			},
		)
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				log := core_logger.FromContext(ctx)
				responseHandler := core_http_response.NewHTTPResponseHandler(
					log, w,
				)

				defer func() {
					if p := recover(); p != nil {
						responseHandler.PanicResponse(
							p,
							"during handler HTTP request got unexpected panic",
						)
					}
				}()

				next.ServeHTTP(w, r)
			},
		)
	}
}
