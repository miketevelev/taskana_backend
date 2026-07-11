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
	areas_postgres_repository "github.com/miketevelev/taskana_backend/internal/features/areas/repository/postgres"
	areas_service "github.com/miketevelev/taskana_backend/internal/features/areas/service"
	areas_transport_http "github.com/miketevelev/taskana_backend/internal/features/areas/transport/http"
	auth_postgres_repository "github.com/miketevelev/taskana_backend/internal/features/auth/repository/postgres"
	auth_service "github.com/miketevelev/taskana_backend/internal/features/auth/service"
	auth_transport_http "github.com/miketevelev/taskana_backend/internal/features/auth/transport/http"
	checklists_postgres_repository "github.com/miketevelev/taskana_backend/internal/features/checklists/repository/postgres"
	checklists_service "github.com/miketevelev/taskana_backend/internal/features/checklists/service"
	checklists_transport_http "github.com/miketevelev/taskana_backend/internal/features/checklists/transport/http"
	heading_postgres_repository "github.com/miketevelev/taskana_backend/internal/features/headings/reporitory/postgres"
	headings_service "github.com/miketevelev/taskana_backend/internal/features/headings/service"
	headings_transport_http "github.com/miketevelev/taskana_backend/internal/features/headings/transport/http"
	projects_postgres_repository "github.com/miketevelev/taskana_backend/internal/features/projects/repository/postgres"
	projects_service "github.com/miketevelev/taskana_backend/internal/features/projects/service"
	projects_transport_http "github.com/miketevelev/taskana_backend/internal/features/projects/transport/http"
	recurring_worker "github.com/miketevelev/taskana_backend/internal/features/recurring/worker"
	task_templates_postgres_repository "github.com/miketevelev/taskana_backend/internal/features/task_templates/reporitory/postgres"
	task_templates_service "github.com/miketevelev/taskana_backend/internal/features/task_templates/service"
	task_templates_transport_http "github.com/miketevelev/taskana_backend/internal/features/task_templates/transport/http"
	tasks_postgres_repository "github.com/miketevelev/taskana_backend/internal/features/tasks/reporitory/postgres"
	tasks_service "github.com/miketevelev/taskana_backend/internal/features/tasks/service"
	tasks_transport_http "github.com/miketevelev/taskana_backend/internal/features/tasks/transport/http"
	user_postgres_repository "github.com/miketevelev/taskana_backend/internal/features/user/repository/postgres"
	user_service "github.com/miketevelev/taskana_backend/internal/features/user/service"
	user_transport_http "github.com/miketevelev/taskana_backend/internal/features/user/transport/http"
	"go.uber.org/zap"

	_ "time/tzdata"
)

func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	// Graceful Shutdown
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	// Logger init
	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("Error initializing logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	// JWT init
	jwtConfig := core_auth.NewJWTConfigMust()
	tokenManager := core_auth.NewTokenManager(jwtConfig)

	// Database init
	pool, err := core_pgx_pool.NewPool(ctx, core_pgx_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	// Init Auth layers (Repository -> Service -> Handler)
	authRepository := auth_postgres_repository.NewAuthRepository(pool)
	authService := auth_service.NewAuthService(authRepository, tokenManager)
	authTransportHTTP := auth_transport_http.NewAuthHTTPHandler(authService)

	// Init User layers (Repository -> Service -> Handler)
	userRepository := user_postgres_repository.NewUserRepository(pool)
	userService := user_service.NewUsersService(userRepository, tokenManager)
	userTransportHTTP := user_transport_http.NewUsersHTTPHandler(
		userService, tokenManager,
	)

	// Init Areas layers (Repository -> Service -> Handler)
	areasRepository := areas_postgres_repository.NewAreasRepository(pool)
	areasService := areas_service.NewAreasService(areasRepository)
	areasTransportHTTP := areas_transport_http.NewAreasHTTPHandler(
		areasService, tokenManager,
	)

	// Init Projects layers (Repository -> Service -> Handler)
	projectsRepository := projects_postgres_repository.NewProjectRepository(pool)
	projectsService := projects_service.NewProjectService(projectsRepository)
	projectsTransportHTTP := projects_transport_http.NewProjectsHTTPHandler(
		projectsService, tokenManager,
	)

	// Init Headings layers (Repository -> Service -> Handler)
	headingsRepository := heading_postgres_repository.NewHeadingRepository(pool)
	headingsService := headings_service.NewHeadingService(headingsRepository)
	headingsTransportHTTP := headings_transport_http.NewHeadingHTTPHandler(
		headingsService, tokenManager,
	)

	// Init Task Templates layers (Repository -> Service -> Handler)
	taskTemplatesRepository := task_templates_postgres_repository.NewTaskTemplateRepository(pool)
	taskTemplatesService := task_templates_service.NewTaskTemplatesService(taskTemplatesRepository)
	taskTemplatesTransportHTTP := task_templates_transport_http.NewTaskTemplatesHTTPHandler(
		taskTemplatesService, tokenManager,
	)

	// Init Tasks layers (Repository -> Service -> Handler)
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTaskService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(
		tasksService, tokenManager,
	)

	// Init Checklists layers (Repository -> Service -> Handler)
	checklistsRepository := checklists_postgres_repository.NewChecklistsRepository(pool)
	checklistsService := checklists_service.NewChecklistsService(checklistsRepository)
	checklistsTransportHTTP := checklists_transport_http.NewChecklistsHTTPHandler(
		checklistsService, tokenManager,
	)

	// Rate Limiter Janitor
	defer authTransportHTTP.Shutdown()

	// Worker
	worker := recurring_worker.NewWorker(tasksService, logger, 24*time.Hour)
	go worker.Run(ctx)

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.CORS(),
		core_http_middleware.RequestId(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	// Api version v1 router
	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(authTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(userTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(areasTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(projectsTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(headingsTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(taskTemplatesTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(tasksTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(checklistsTransportHTTP.Routes()...)

	httpServer.RegisterAPIRoutes(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
