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

	CreateProject(
		ctx context.Context,
		userID uuid.UUID,
		project domain.Project,
	) (domain.Project, error)
}

func NewProjectService(
	projectsRepository ProjectsRepository,
) *ProjectService {
	return &ProjectService{
		projectsRepository: projectsRepository,
	}
}
