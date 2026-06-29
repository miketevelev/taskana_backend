package core_http_middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_http_response "github.com/miketevelev/taskana_backend/internal/core/transport/http/response"
)

type rateLimiter struct {
	mu      sync.Mutex
	entries map[string][]time.Time
	limit   int
	window  time.Duration
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		entries: make(map[string][]time.Time),
		limit:   limit,
		window:  window,
	}
}

func (rl *rateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	times := rl.entries[key]
	filtered := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) >= rl.limit {
		rl.entries[key] = filtered
		return false
	}

	rl.entries[key] = append(filtered, now)
	return true
}

func AuthRateLimit(limit int, window time.Duration) Middleware {
	limiter := newRateLimiter(limit, window)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				log := core_logger.FromContext(ctx)
				responseHandler := core_http_response.NewHTTPResponseHandler(
					log, w,
				)

				key := clientKey(r)
				if !limiter.allow(key) {
					responseHandler.ErrorResponse(
						core_errors.ErrTooManyRequests,
						"rate limit exceeded",
					)
					return
				}

				next.ServeHTTP(w, r)
			},
		)
	}
}

func clientKey(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil {
			ip = host
		} else {
			ip = r.RemoteAddr
		}
	}
	if idx := strings.Index(ip, ","); idx >= 0 {
		ip = strings.TrimSpace(ip[:idx])
	}

	email := r.Header.Get("X-Auth-Email")
	if email != "" {
		return ip + ":" + strings.ToLower(email)
	}

	return ip
}
