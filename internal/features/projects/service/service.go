package projects_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/miketevelev/taskana_backend/internal/core/domain/project"
)

type ProjectService struct {
	projectsRepository ProjectsRepository
}

type ProjectsRepository interface {
	GetProject(
		ctx context.Context,
		userID uuid.UUID,
		projectID uuid.UUID,
	) (domain_project.Project, error)

	GetProjects(
		ctx context.Context,
		userID uuid.UUID,
		limit *int,
		offset *int,
	) ([]domain_project.Project, error)

	CreateProject(
		ctx context.Context,
		userID uuid.UUID,
		project domain_project.Project,
	) (domain_project.Project, error)

	ChangePosition(
		ctx context.Context,
		userID uuid.UUID,
		project domain_project.Project,
		oldPosition int,
	) (domain_project.Project, error)

	PatchProject(
		ctx context.Context,
		userID uuid.UUID,
		project domain_project.Project,
	) (domain_project.Project, error)

	PatchProjectWithAreaChange(
		ctx context.Context,
		userID uuid.UUID,
		project domain_project.Project,
		oldPosition int,
		oldAreaID *uuid.UUID,
	) (domain_project.Project, error)

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
