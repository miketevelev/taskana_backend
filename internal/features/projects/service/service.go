package projects_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain"
)

type ProjectService struct {
	projectsRepository ProjectsRepository
}

type ProjectsRepository interface {
	GetProject(
		ctx context.Context,
		userID uuid.UUID,
		projectID uuid.UUID,
	) (domain.Project, error)

	GetProjects(
		ctx context.Context,
		userID uuid.UUID,
		limit *int,
		offset *int,
	) ([]domain.Project, error)

	CreateProject(
		ctx context.Context,
		userID uuid.UUID,
		project domain.Project,
	) (domain.Project, error)

	PatchProject(
		ctx context.Context,
		userID uuid.UUID,
		project domain.Project,
	) (domain.Project, error)

	PatchProjectWithReordering(
		ctx context.Context,
		userID uuid.UUID,
		project domain.Project,
		oldPosition int,
		oldAreaID *uuid.UUID,
	) (domain.Project, error)

	DeleteProject(
		ctx context.Context,
		userID uuid.UUID,
		projectID uuid.UUID,
	) error
}

func NewProjectService(
	projectsRepository ProjectsRepository,
) *ProjectService {
	return &ProjectService{
		projectsRepository: projectsRepository,
	}
}
