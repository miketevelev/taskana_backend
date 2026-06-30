package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_auth "github.com/miketevelev/taskana_backend/internal/core/auth"
	core_config "github.com/miketevelev/taskana_backend/internal/core/config"
	core_logger "github.com/miketevelev/taskana_backend/internal/core/logger"
	core_pgx_pool "github.com/miketevelev/taskana_backend/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/miketevelev/taskana_backend/internal/core/transport/http/middleware"
	core_http_server "github.com/miketevelev/taskana_backend/internal/core/transport/http/server"
	auth_postgres_repository "github.com/miketevelev/taskana_backend/internal/features/auth/repository/postgres"
	auth_service "github.com/miketevelev/taskana_backend/internal/features/auth/service"
	auth_transport_http "github.com/miketevelev/taskana_backend/internal/features/auth/transport/http"
	"go.uber.org/zap"

	_ "time/tzdata"
)

func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("Error initializing logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	jwtConfig := core_auth.NewJWTConfigMust()
	tokenManager := core_auth.NewTokenManager(jwtConfig)

	pool, err := core_pgx_pool.NewPool(ctx, core_pgx_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	authRepository := auth_postgres_repository.NewAuthRepository(pool)
	authService := auth_service.NewAuthService(authRepository, tokenManager)
	authTransportHTTP := auth_transport_http.NewAuthHTTPHandler(authService)

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.CORS(),
		core_http_middleware.RequestId(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(authTransportHTTP.Routes()...)

	httpServer.RegisterAPIRoutes(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
